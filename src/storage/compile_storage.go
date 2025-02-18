package storage

import (
	"errors"
	"nullableocean-postupashki/src/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage interface {
	SaveTask(task *models.CompileTask) error
	FindTask(id string) (*models.CompileTask, error)
}

type RamStorage struct {
	tasksMap map[string]*models.CompileTask
}

func NewRamStorage() Storage {
	return &RamStorage{
		tasksMap: map[string]*models.CompileTask{},
	}
}

func (rs *RamStorage) SaveTask(task *models.CompileTask) error {
	if _, exist := rs.tasksMap[task.Uuid]; exist {
		return nil
	}

	rs.tasksMap[task.Uuid] = task
	return nil
}

func (rs *RamStorage) FindTask(id string) (*models.CompileTask, error) {
	t, exist := rs.tasksMap[id]
	if !exist {
		return nil, ErrTaskNotFound
	}

	return t, nil
}
