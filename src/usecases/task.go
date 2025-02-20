package usecases

import "nullableocean-postupashki/src/domain"

type Task interface {
	Get(id string) (*domain.Task, error)
	CheckStatus(*domain.Task) domain.TaskStatus
	CreateAndProcess(*domain.Task) (*domain.Task, error)
}
