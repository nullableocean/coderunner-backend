package processor

import (
	"codeproccesor/src/domain"
	"codeproccesor/src/usecases"
	"codeproccesor/src/usecases/service/coderunner"
	"context"
	"fmt"
	"time"
)

const (
	timeLimit = 60 * time.Second
)

type CodeProcessor struct {
	commiter usecases.Commiter
	runner   *coderunner.CodeRunner
}

func NewCodeProcessor(commiter usecases.Commiter, runner *coderunner.CodeRunner) usecases.Processor {
	return &CodeProcessor{
		commiter: commiter,
		runner:   runner,
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
		result.Error = err.Error()
		result.IsSuccess = false

		if err := proc.saveResult(result); err != nil {
			return err
		}

		return err
	}

	result.Output = string(execInfo.Output)
	result.IsSuccess = execInfo.IsSuccess()

	err = proc.saveResult(result)
	if err != nil {
		return fmt.Errorf("save result error: %s", err)
	}

	return nil
}

func (proc *CodeProcessor) saveResult(result domain.Result) error {
	return proc.commiter.Commit(result)
}
