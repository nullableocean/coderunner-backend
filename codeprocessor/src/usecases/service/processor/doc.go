package processor

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const (
	runnerDockerfileCtx = "./runnerdocker"
	runnerDockerfile    = "runner.Dockerfile"
	imageName           = "runner:latest"
	runnerUser          = "1001:1001"
	runnerFilesDir      = "/app/"
)

type Task struct {
	Code     string
	Compiler string
}

type ExecResponse struct {
	Out       string
	IsSuccess bool
	ExitCode  int
}

var (
	testTask = Task{
		Code: `#include <stdio.h>
int main() {
    printf("Hello from sandbox!\\n");
    return 0;
}`,
		Compiler: "clang",
	}
)

func doc() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	l, _ := cli.ImageList(ctx, image.ListOptions{})
	imageId := ""
	for _, i := range l {
		for _, t := range i.RepoTags {
			fmt.Println(t)
			if t == imageName {
				imageId = i.ID
				break
			}
		}
		if imageId != "" {
			break
		}
	}
	if imageId != "" {
		cli.ImageRemove(ctx, imageId, image.RemoveOptions{})
	}

	// CREATE IMAGE
	fmt.Println("Собираем раннер-образ...")

	tarBuildCtx, err := archive.TarWithOptions(runnerDockerfileCtx, &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ImageBuild(ctx, tarBuildCtx, types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: runnerDockerfile,
		NoCache:    true,
	})
	if err != nil {
		panic(err)
	}
	io.ReadAll(resp.Body) //need for build (magick)
	resp.Body.Close()

	// CONTAINER UP
	fmt.Println("Создаем раннер-контейнер...")

	config := container.Config{
		User:         runnerUser,
		AttachStdout: true,
		AttachStderr: true,
		Image:        imageName,
		Cmd:          []string{"tail", "-f", "/dev/null"},
	}

	hostConfig := container.HostConfig{
		AutoRemove: true,
		Resources: container.Resources{
			Memory:   512 * 1024 * 1024, // Ограничение RAM: 512MB
			CPUQuota: 50000,             // Лимит CPU: 50% (в единицах cgroups)
		},
		NetworkMode: "none",
		CapDrop:     []string{"ALL"},
	}

	respContainer, err := cli.ContainerCreate(ctx, &config, &hostConfig, nil, nil, "")
	if err != nil {
		panic(err)
	}
	defer func() {
		removeOpts := container.RemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		}
		if err := cli.ContainerRemove(ctx, respContainer.ID, removeOpts); err != nil {
			log.Printf("Warning: container remove failed: %v", err)
		}
	}()
	// defer cli.ContainerRemove(ctx, respContainer.ID, container.RemoveOptions{})

	cid := respContainer.ID

	// COPY IN CONTAINER
	fmt.Println("Создаем файл с кодом...")

	codePath := runnerFilesDir + getFileName(&testTask)
	err = CopyToContainer(ctx, cli, cid, codePath, &testTask)
	if err != nil {
		panic(err)
	}

	// COPY IN CONTAINER
	fmt.Println("Поднимаем контейнер...")
	err = cli.ContainerStart(ctx, cid, container.StartOptions{})
	if err != nil {
		panic(err)
	}

	inspect, err := cli.ContainerInspect(ctx, cid)
	if err != nil {
		panic(err)
	}
	if !inspect.State.Running {
		panic("Контейнер не запущен!")
	}

	// COMPILE
	fmt.Println("Компилируем...")

	buildPath := runnerFilesDir + "build"
	compileCmd := []string{"gcc", "-o", buildPath, codePath}

	compileRes := execCommand(ctx, cli, cid, compileCmd)

	lsres := execCommand(ctx, cli, cid, []string{"ls", "-l"})
	fmt.Printf("LS:\n%s\nУспешно: %t\nКод: %d\n\n",
		lsres.Out, lsres.IsSuccess, lsres.ExitCode)

	pwdres := execCommand(ctx, cli, cid, []string{"pwd"})
	fmt.Printf("LS:\n%s\nУспешно: %t\nКод: %d\n\n",
		pwdres.Out, pwdres.IsSuccess, pwdres.ExitCode)

	// EXECUTE
	fmt.Println("Запускаем...")
	execRes := execCommand(ctx, cli, cid, []string{buildPath})

	fmt.Printf("Компиляция:\n%s\nУспешно: %t\nКод: %d\n\n",
		compileRes.Out, compileRes.IsSuccess, compileRes.ExitCode)

	fmt.Printf("Выполнение:\n%s\nУспешно: %t\nКод: %d\n",
		execRes.Out, execRes.IsSuccess, execRes.ExitCode)

}

func CopyToContainer(ctx context.Context, cli *client.Client, cid string, path string, task *Task) error {
	// Проверяем состояние контейнера
	inspect, err := cli.ContainerInspect(ctx, cid)
	if err != nil {
		return err
	}

	if !inspect.State.Running && inspect.State.Status != "created" {
		return fmt.Errorf("container is in invalid state: %s", inspect.State.Status)
	}

	tarBuff := bytes.Buffer{}
	tw := tar.NewWriter(&tarBuff)
	thdr := tar.Header{
		Name: filepath.Base(path),
		Size: int64(len(task.Code)),
		Mode: 0644,
	}
	tw.WriteHeader(&thdr)
	tw.Write([]byte(task.Code))
	tw.Close()

	err = cli.CopyToContainer(ctx, cid, filepath.Dir(path), &tarBuff, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                true,
	})
	if err != nil {
		return err
	}

	return nil
}

func getFileName(t *Task) string {
	return "code.c"
}

func execCommand(ctx context.Context, cli *client.Client, cid string, cmd []string) ExecResponse {

	execOpts := container.ExecOptions{
		User:         "1001", // or username
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}

	execId, err := cli.ContainerExecCreate(ctx, cid, execOpts)
	if err != nil {
		return ExecResponse{
			Out:       fmt.Sprintf("Exec create failed: %v", err),
			IsSuccess: false,
			ExitCode:  -1,
		}
	}

	resp, err := cli.ContainerExecAttach(ctx, execId.ID, container.ExecAttachOptions{})
	if err != nil {
		return ExecResponse{
			Out:       fmt.Sprintf("Exec attach failed: %v", err),
			IsSuccess: false,
			ExitCode:  -1,
		}
	}
	defer resp.Close()

	outBuff := bytes.Buffer{}
	_, err = outBuff.ReadFrom(resp.Reader)
	if err != nil {
		return ExecResponse{
			Out:       fmt.Sprintf("Read stdout/stderr failed: %v", err),
			IsSuccess: false,
			ExitCode:  -1,
		}
	}

	execInspect, err := cli.ContainerExecInspect(ctx, execId.ID)
	if err != nil {
		return ExecResponse{
			Out:       fmt.Sprintf("Exec inspect failed: %v", err),
			IsSuccess: false,
			ExitCode:  -1,
		}
	}

	return ExecResponse{
		Out:       outBuff.String(),
		IsSuccess: execInspect.ExitCode == 0,
		ExitCode:  execInspect.ExitCode,
	}
}

func safeExec(ctx context.Context, cli *client.Client, containerID string, cmd []string, timeout time.Duration) (int, string, error) {
	// Создаем контекст с таймаутом
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	execID, err := cli.ContainerExecCreate(timeoutCtx, containerID, container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		User:         "1001",
	})
	if err != nil {
		return -1, "", fmt.Errorf("exec create failed: %v", err)
	}

	// Захватываем вывод
	output := new(bytes.Buffer)
	resp, err := cli.ContainerExecAttach(timeoutCtx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return -1, "", fmt.Errorf("exec attach failed: %v", err)
	}
	defer resp.Close()

	// Читаем вывод асинхронно
	go func() {
		io.Copy(output, resp.Reader)
	}()

	// Ждем завершения
	for {
		select {
		case <-timeoutCtx.Done():
			return -1, "", fmt.Errorf("execution timeout")
		default:
			inspect, err := cli.ContainerExecInspect(timeoutCtx, execID.ID)
			if err != nil {
				return -1, "", fmt.Errorf("inspect failed: %v", err)
			}
			if !inspect.Running {
				return inspect.ExitCode, output.String(), nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
