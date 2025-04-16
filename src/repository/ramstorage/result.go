package ramstorage

import (
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
)

type ResultRepository struct {
	rMap map[string]*domain.Result
}

func NewResultRepository() repository.ResultRepository {
	return &ResultRepository{
		rMap: map[string]*domain.Result{},
	}
}

func (rs *ResultRepository) Post(res *domain.Result) error {
	rs.rMap[res.TaskUuid] = res

	return nil
}

func (rs *ResultRepository) GetByTaskId(id string) (*domain.Result, error) {
	res, exist := rs.rMap[id]
	if !exist {
		return nil, repository.ErrNotFound
	}

	return res, nil
}
