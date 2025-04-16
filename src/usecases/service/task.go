package service

import (
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"nullableocean-postupashki/src/usecases"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo   repository.TaskRepository
	taskSender repository.TaskSender
}

func NewTaskService(repository repository.TaskRepository, sender repository.TaskSender) usecases.Task {
	return &TaskService{
		taskRepo:   repository,
		taskSender: sender,
	}
}

func (s *TaskService) Get(id string) (*domain.Task, error) {
	return s.taskRepo.GetById(id)
}

func (s *TaskService) CheckStatus(task *domain.Task) domain.TaskStatus {
	return task.Status
}

func (s *TaskService) UpdateStatus(task *domain.Task, newStatus domain.TaskStatus) error {
	task.Status = newStatus
	return s.taskRepo.Update(task)
}

func (s *TaskService) CreateAndProcess(data *domain.Task) (*domain.Task, error) {
	task, err := s.create(data)
	if err == nil {
		s.taskSender.Send(task)
	}

	return task, err
}

func (s *TaskService) create(task *domain.Task) (*domain.Task, error) {
	task.Uuid = uuid.NewString()
	task.Status = domain.InProgress

	err := s.taskRepo.Post(task)
	return task, err
}

func (s *TaskService) process(task *domain.Task) {
	go func(task *domain.Task) {
		time.Sleep(20 * time.Second)
		task.Status = domain.Ready
	}(task)
}
