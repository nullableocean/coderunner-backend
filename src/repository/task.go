package repository

import (
	"nullableocean-postupashki/src/domain"
)

type TaskRepository interface {
	Post(*domain.Task) error
	GetById(id string) (*domain.Task, error)
}
