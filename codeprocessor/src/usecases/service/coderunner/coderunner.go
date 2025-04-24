package coderunner

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type CodeRunner struct {
	client *client.Client
	logger *log.Logger

	runnerDockerfilePath string
}

func NewCodeRunner(logger *log.Logger, dockerClient *client.Client, runnerDockerfilePath string) (*CodeRunner, error) {
	runner := &CodeRunner{
		client:               dockerClient,
		logger:               logger,
		runnerDockerfilePath: runnerDockerfilePath,
	}

	err := runner.init()

	return runner, err
}

func (r *CodeRunner) init() error {
	r.logger.Printf("build image...")

	err := r.buildImage()
	if err != nil {
		return err
	}

	tries := 5
	ctx := context.Background()
	imageExist := false
	for tries != 0 {
		time.Sleep(5 * time.Second)

		imageExist, err = r.imageExists(ctx)
		if err != nil {
			return err
		}
		if imageExist {
			r.logger.Printf("runner image ready")
			break
		}

		tries--
	}

	if !imageExist {
		return errors.New("runner image not exist")
	}

	return nil
}

func (r *CodeRunner) Execute(ctx context.Context, compiler Compiler, code string) (*ExecInfo, error) {
	if compiler == UNKNOWN {
		return &ExecInfo{
			Output:   []byte("unknown compiler"),
			ExitCode: -1,
		}, nil
	}

	c, err := r.runContainer(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		go r.remove(context.Background(), c)
	}()

	err = r.createSourceFile(ctx, c, compiler, code)
	if err != nil {
		return nil, err
	}

	compExInfo, err := r.compileSourseFile(ctx, c, compiler)
	if err != nil {
		return nil, err
	}
	if !compExInfo.IsSuccess() {
		return compExInfo, nil
	}

	return r.runCode(ctx, c, compiler)
}

func (r *CodeRunner) runContainer(ctx context.Context) (*Container, error) {
	config := container.Config{
		User:         runnerUser,
		AttachStdout: true,
		AttachStderr: true,
		Image:        runnerImage,
		Cmd:          []string{"tail", "-f", "/dev/null"},
	}

	hostConfig := container.HostConfig{
		AutoRemove: true,
		Resources: container.Resources{
			Memory:   100 * 1024 * 1024, // 100MB
			CPUQuota: 50000,             // CPU: 50% (units cgroups)
		},
		NetworkMode: "none",
		CapDrop:     []string{"ALL"},
	}

	resp, err := r.client.ContainerCreate(ctx, &config, &hostConfig, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("container create error: %s", err)
	}

	c := &Container{
		ID: resp.ID,
	}

	err = r.client.ContainerStart(ctx, c.ID, container.StartOptions{})
	if err != nil {
		return nil, fmt.Errorf("container start error: %s", err)
	}

	return c, nil
}

func (r *CodeRunner) isRunning(ctx context.Context, c *Container) (bool, error) {
	inspect, err := r.client.ContainerInspect(ctx, c.ID)
	if err != nil {
		return false, fmt.Errorf("container inspect error: %s", err)
	}

	return inspect.State.Running, nil
}

func (r *CodeRunner) createSourceFile(ctx context.Context, c *Container, compiler Compiler, code string) error {
	tarBuff := &bytes.Buffer{}
	tarWriter := tar.NewWriter(tarBuff)
	tarHdr := tar.Header{
		Name: compiler.SourseFileName(),
		Size: int64(len(code)),
		Mode: 0644,
	}
	err := tarWriter.WriteHeader(&tarHdr)
	if err != nil {
		return fmt.Errorf("tar header error: %s", err)
	}

	_, err = tarWriter.Write([]byte(code))
	if err != nil {
		return fmt.Errorf("tar write error: %s", err)
	}

	err = r.client.CopyToContainer(ctx, c.ID, runnerWorkdir, tarBuff, container.CopyToContainerOptions{})
	if err != nil {
		return fmt.Errorf("container code copy error: %s", err)
	}

	return nil
}

func (r *CodeRunner) compileSourseFile(ctx context.Context, c *Container, compiler Compiler) (*ExecInfo, error) {
	if compiler.IsInterpreter() {
		return &ExecInfo{
			ExitCode: 0,
		}, nil
	}

	exInfo, err := r.execCmd(ctx, c, compiler.CompileCommand())
	if err != nil {
		return nil, fmt.Errorf("container compile file error: %s", err)
	}

	return exInfo, nil
}

func (r *CodeRunner) runCode(ctx context.Context, c *Container, compiler Compiler) (*ExecInfo, error) {
	r.logger.Printf("run code. compiler %s\n", compiler.String())

	exInfo, err := r.execCmd(ctx, c, compiler.RunCommand())
	if err != nil {
		return nil, fmt.Errorf("container execute code error: %s", err)
	}

	r.logger.Printf("run code. exit code %d\n", exInfo.ExitCode)

	return exInfo, nil
}

func (r *CodeRunner) execCmd(ctx context.Context, c *Container, cmd []string) (*ExecInfo, error) {
	execInfo := &ExecInfo{
		ExitCode: -1,
	}

	execId, err := r.client.ContainerExecCreate(ctx, c.ID, container.ExecOptions{
		User:         runnerUser,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, fmt.Errorf("container exec error: %s", err)
	}

	hijRes, err := r.client.ContainerExecAttach(ctx, execId.ID, container.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("container attach to exec error: %s", err)
	}
	defer hijRes.Close()

	out, err := io.ReadAll(hijRes.Reader)
	if err != nil {
		return nil, fmt.Errorf("container read exec out error: %s", err)
	}

	inspect, err := r.client.ContainerExecInspect(ctx, execId.ID)
	if err != nil {
		return nil, fmt.Errorf("container inspect exec error: %s", err)
	}

	execInfo.Output = out
	execInfo.ExitCode = inspect.ExitCode

	return execInfo, nil
}

func (r *CodeRunner) buildImage() error {
	ctx := context.Background()

	buildCtx, err := archive.TarWithOptions(r.runnerDockerfilePath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("build context error: %s", err)
	}

	resp, err := r.client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{runnerImage},
		Dockerfile: "Dockerfile",
		NoCache:    true,
	})
	if err != nil {
		return fmt.Errorf("build image error: %s", err)
	}

	defer resp.Body.Close()
	output, _ := io.ReadAll(resp.Body)
	r.logger.Printf("image build output: %s\n\n", output)

	return nil
}

func (r *CodeRunner) remove(ctx context.Context, c *Container) error {
	err := r.client.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("container remove error: %s", err)
	}

	return nil
}

func (r *CodeRunner) imageExists(ctx context.Context) (bool, error) {
	images, err := r.client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("list images error: %w", err)
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == runnerImage {
				return true, nil
			}
		}
	}

	return false, nil
}
