package repository

import (
	"nullableocean-postupashki/src/domain"
)

type UserRepository interface {
	Create(*domain.User) (*domain.User, error)
	GetByLogin(login string) (*domain.User, error)
	GetById(id int64) (*domain.User, error)
}
