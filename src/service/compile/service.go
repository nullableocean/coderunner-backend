package compile

import (
	"nullableocean-postupashki/src/models"
	"nullableocean-postupashki/src/storage"
	"time"

	"github.com/google/uuid"
)

type CompileTaskService struct {
	storage storage.CompileStorage
}

func NewCompileService(storage storage.CompileStorage) *CompileTaskService {
	return &CompileTaskService{
		storage: storage,
	}
}

type TaskCreateData struct {
	Code         string
	CompilerName string
}

func (s *CompileTaskService) CreateTask(data *TaskCreateData) (*models.CompileTask, error) {
	task := &models.CompileTask{
		Uuid:         uuid.NewString(),
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

func (s *CompileTaskService) GetTask(id string) (*models.CompileTask, error) {
	return s.storage.FindTask(id)
}

func (s *CompileTaskService) CheckTaskStatus(task *models.CompileTask) models.TaskStatus {
	return task.Status
}

func (s *CompileTaskService) handleTask(task *models.CompileTask) {
	go func(task *models.CompileTask) {
		time.Sleep(20 * time.Second)
		task.Status = models.Ready
	}(task)
}
