package usecases

import "nullableocean-postupashki/src/domain"

type Session interface {
	Start(*domain.User) (*domain.Session, error)
	Get(sid string) (*domain.Session, error)
	Destroy(sid string) error
	DestroyAllByUser(userId int64) error
}
