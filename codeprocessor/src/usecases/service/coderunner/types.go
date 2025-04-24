package coderunner

const (
	runnerImage             = "runner:latest"
	runnerContextDockerfile = "Dockerfile"
	runnerUser              = "1001"
	runnerWorkdir           = "/app/"
	buildFileName           = "a.out"
)

type Container struct {
	ID string
}

type ExecInfo struct {
	Output   []byte
	ExitCode int
}

func (i ExecInfo) IsSuccess() bool {
	return i.ExitCode == 0
}

func runnerFilePath(file string) string {
	return runnerWorkdir + file
}
