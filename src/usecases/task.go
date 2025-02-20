package usecases

import "nullableocean-postupashki/src/domain"

type TaskCreateData struct {
	Code         string
	CompilerName string
}

type TaskService interface {
	GetTask(id string) (*domain.Task, error)
	CheckStatus(task *domain.Task) domain.TaskStatus
	CreateAndProcess(data *TaskCreateData) (*domain.Task, error)
}
