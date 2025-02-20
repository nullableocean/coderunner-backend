package ramstorage

import (
	"nullableocean-postupashki/src/models"
	"nullableocean-postupashki/src/storage"
)

type RamStorage struct {
	tasksMap map[string]*models.CompileTask
}

func NewRamStorage() storage.CompileStorage {
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
		return nil, storage.ErrNotFound
	}

	return t, nil
}
