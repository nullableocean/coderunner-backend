package processor

import (
	"codeproccesor/src/domain"
	"codeproccesor/src/usecases"
	"codeproccesor/src/usecases/service/coderunner"
	"context"
	"fmt"
	"log"
	"time"
)

const (
	timeLimit = 60 * time.Second
)

type CodeProcessor struct {
	commiter usecases.Commiter
	runner   *coderunner.CodeRunner
	logger   *log.Logger
}

func NewCodeProcessor(logger *log.Logger, commiter usecases.Commiter, runner *coderunner.CodeRunner) usecases.Processor {
	return &CodeProcessor{
		commiter: commiter,
		runner:   runner,
		logger:   logger,
	}
}

func (proc *CodeProcessor) Process(task domain.Task) error {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(timeLimit))

	compiler := coderunner.CompilerFromString(task.CompilerName)

	result := domain.Result{
		TaskUuid:  task.Uuid,
		Output:    "",
		Error:     "",
		IsSuccess: false,
	}

	execInfo, err := proc.runner.Execute(ctx, compiler, task.Code)
	if err != nil {
		proc.logger.Printf("process execute error: %s\n", err)

		result.Error = err.Error()
		result.IsSuccess = false
	} else {
		result.Output = string(execInfo.Output)
		result.IsSuccess = execInfo.IsSuccess()
	}

	if err := proc.saveResult(result); err != nil {
		return fmt.Errorf("save result error: %s", err)
	}

	return err
}

func (proc *CodeProcessor) saveResult(result domain.Result) error {
	proc.logger.Printf("commit result. output: %.20s", result.Output)

	return proc.commiter.Commit(result)
}
