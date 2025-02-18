package service

import (
	"nullableocean-postupashki/src/models"
	"nullableocean-postupashki/src/storage"
	"nullableocean-postupashki/src/utils"
	"time"
)

type CompileTaskManager struct {
	storage storage.Storage
}

func NewCompileManager(storage storage.Storage) *CompileTaskManager {
	return &CompileTaskManager{
		storage: storage,
	}
}

type TaskCreateData struct {
	Code         string `json:"code"`
	CompilerName string `json:"compiler_name"`
}

func (s *CompileTaskManager) CreateTask(data *TaskCreateData) (*models.CompileTask, error) {
	task := &models.CompileTask{
		Uuid:         utils.CreateUuid(),
		Code:         data.Code,
		CompilerName: data.CompilerName,
		Status:       models.InProgress,
	}

	err := s.storage.SaveTask(task)
	if err == nil {
		s.handleTask(task)
	}

	return task, err
}

func (s *CompileTaskManager) GetTask(id string) (*models.CompileTask, error) {
	return s.storage.FindTask(id)
}

func (s *CompileTaskManager) CheckTaskStatus(task *models.CompileTask) models.TaskStatus {
	return task.Status
}

func (s *CompileTaskManager) handleTask(task *models.CompileTask) {
	go func(task *models.CompileTask) {
		time.Sleep(20 * time.Second)
		task.Status = models.Ready
	}(task)
}
