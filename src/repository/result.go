package repository

import "nullableocean-postupashki/src/domain"

type ResultRepository interface {
	Post(*domain.Result) error
	GetByTaskId(taskId string) (*domain.Result, error)
}
