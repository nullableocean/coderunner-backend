package service

import (
	"fmt"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"nullableocean-postupashki/src/usecases"
)

type ResultService struct {
	repo        repository.ResultRepository
	taskService usecases.Task
}

func NewResultService(repo repository.ResultRepository, taskService usecases.Task) usecases.Result {
	return &ResultService{
		repo:        repo,
		taskService: taskService,
	}
}

func (rs *ResultService) SaveResult(result *domain.Result) error {
	t, err := rs.taskService.Get(result.TaskUuid)
	if err != nil {
		return fmt.Errorf("save result error: cannot get relation task error: %s", err)
	} else if t == nil {
		return fmt.Errorf("save result error: cannot get relation task with id: %s", result.TaskUuid)
	}

	res, _ := rs.repo.GetByTaskId(result.TaskUuid)
	if res != nil {
		return nil
	}

	err = rs.repo.Post(result)

	if err == nil {
		err = rs.taskService.UpdateStatus(t, domain.Ready)
	}

	return err
}

func (rs *ResultService) GetTaskResult(task *domain.Task) (*domain.Result, error) {
	if task.Uuid == "" {
		return nil, fmt.Errorf("get result for task error: task uuid is empty")
	}

	return rs.repo.GetByTaskId(task.Uuid)
}
