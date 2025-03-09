package ramstorage

import (
	"errors"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"sync/atomic"
)

type UserRepository struct {
	idCounter     int64
	usersLoginMap map[string]*domain.User
	usersIdMap    map[int64]*domain.User
}

func NewUserRepository() repository.UserRepository {
	return &UserRepository{
		idCounter:     0,
		usersLoginMap: map[string]*domain.User{},
		usersIdMap:    map[int64]*domain.User{},
	}
}

func (r *UserRepository) Create(user *domain.User) (*domain.User, error) {
	if user.HashedPassword == "" {
		return nil, errors.New("password hash empty")
	}

	if _, exist := r.usersLoginMap[user.Login]; exist {
		return nil, errors.New("user already register")
	}

	user.Id = atomic.AddInt64(&r.idCounter, 1)
	r.usersLoginMap[user.Login] = user
	r.usersIdMap[user.Id] = user

	return user, nil
}

func (r *UserRepository) GetByLogin(login string) (*domain.User, error) {
	user, exist := r.usersLoginMap[login]
	if !exist {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) GetById(id int64) (*domain.User, error) {
	user, exist := r.usersIdMap[id]
	if !exist {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}
