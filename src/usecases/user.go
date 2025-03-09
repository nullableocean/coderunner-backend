package usecases

import "nullableocean-postupashki/src/domain"

type User interface {
	Create(domain.User) (*domain.User, error)
	SignIn(*domain.User) (*domain.Session, error)
	CreateSession(*domain.User) (*domain.Session, error)
}
