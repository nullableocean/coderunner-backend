package service

import (
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"time"

	"github.com/google/uuid"
)

type CompileTaskService struct {
	repo repository.TaskRepository
}

func NewCompileTaskService(repository repository.TaskRepository) *CompileTaskService {
	return &CompileTaskService{
		repo: repository,
	}
}

func (s *CompileTaskService) GetTask(id string) (*domain.Task, error) {
	return s.repo.GetById(id)
}

func (s *CompileTaskService) CheckStatus(task *domain.Task) domain.TaskStatus {
	return task.Status
}

type TaskCreateData struct {
	Code         string
	CompilerName string
}

func (s *CompileTaskService) CreateAndProcess(data *TaskCreateData) (*domain.Task, error) {
	task, err := s.create(data)
	if err == nil {
		s.process(task)
	}

	return task, err
}

func (s *CompileTaskService) create(data *TaskCreateData) (*domain.Task, error) {
	task := &domain.Task{
		Uuid:         uuid.NewString(),
		Code:         data.Code,
		CompilerName: data.CompilerName,
		Status:       domain.InProgress,
	}

	err := s.repo.Post(task)
	return task, err
}

func (s *CompileTaskService) process(task *domain.Task) {
	go func(task *domain.Task) {
		time.Sleep(20 * time.Second)
		task.Status = domain.Ready
	}(task)
}
