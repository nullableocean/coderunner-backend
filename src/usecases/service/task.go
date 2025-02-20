package service

import (
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"nullableocean-postupashki/src/usecases"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	repo repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) usecases.Task {
	return &Task{
		repo: repository,
	}
}

func (s *Task) Get(id string) (*domain.Task, error) {
	return s.repo.GetById(id)
}

func (s *Task) CheckStatus(task *domain.Task) domain.TaskStatus {
	return task.Status
}

func (s *Task) CreateAndProcess(data *domain.Task) (*domain.Task, error) {
	task, err := s.create(data)
	if err == nil {
		s.process(task)
	}

	return task, err
}

func (s *Task) create(task *domain.Task) (*domain.Task, error) {
	task.Uuid = uuid.NewString()
	task.Status = domain.InProgress

	err := s.repo.Post(task)
	return task, err
}

func (s *Task) process(task *domain.Task) {
	go func(task *domain.Task) {
		time.Sleep(20 * time.Second)
		task.Status = domain.Ready
	}(task)
}
