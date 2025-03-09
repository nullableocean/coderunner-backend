package repository

import "nullableocean-postupashki/src/domain"

type Session interface {
	Save(*domain.Session) error
	Delete(sid string) error
	Get(sid string) (*domain.Session, error)
	GetByUserId(id int64) ([]*domain.Session, error)
}
