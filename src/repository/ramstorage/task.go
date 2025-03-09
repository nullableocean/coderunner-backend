package ramstorage

import (
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
)

type TaskRepository struct {
	tasksMap map[string]*domain.Task
}

func NewTaskRepository() repository.TaskRepository {
	return &TaskRepository{
		tasksMap: map[string]*domain.Task{},
	}
}

func (rs *TaskRepository) Post(task *domain.Task) error {
	if _, exist := rs.tasksMap[task.Uuid]; exist {
		return nil
	}

	rs.tasksMap[task.Uuid] = task
	return nil
}

func (rs *TaskRepository) GetById(id string) (*domain.Task, error) {
	t, exist := rs.tasksMap[id]
	if !exist {
		return nil, repository.ErrTaskNotFound
	}

	return t, nil
}
