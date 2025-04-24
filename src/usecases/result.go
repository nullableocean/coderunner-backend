package usecases

import "nullableocean-postupashki/src/domain"

type Result interface {
	SaveResult(result *domain.Result) error
	GetTaskResult(task *domain.Task) (*domain.Result, error)
}
