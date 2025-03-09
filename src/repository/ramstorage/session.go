package ramstorage

import (
	"maps"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"slices"
)

type SessionRepository struct {
	sessionMap       map[string]*domain.Session
	sessionUserIdMap map[int64]map[string]*domain.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessionMap:       make(map[string]*domain.Session),
		sessionUserIdMap: make(map[int64]map[string]*domain.Session),
	}
}

func (repo *SessionRepository) Save(session *domain.Session) error {
	if ses, _ := repo.Get(session.SessionId); ses != nil {
		return repository.ErrSessionAlreadyExist
	}

	repo.sessionMap[session.SessionId] = session

	if _, exist := repo.sessionUserIdMap[session.UserId]; !exist {
		repo.sessionUserIdMap[session.UserId] = make(map[string]*domain.Session)
	}

	repo.sessionUserIdMap[session.UserId][session.SessionId] = session

	return nil
}

func (repo *SessionRepository) Get(sid string) (*domain.Session, error) {
	session, exist := repo.sessionMap[sid]
	if !exist {
		return nil, repository.ErrSessionNotFound
	}

	return session, nil
}

func (repo *SessionRepository) Delete(sid string) error {
	session, err := repo.Get(sid)
	if err != nil {
		return err
	}

	delete(repo.sessionMap, sid)
	delete(repo.sessionUserIdMap[session.UserId], sid)

	return nil
}

func (repo *SessionRepository) GetByUserId(id int64) ([]*domain.Session, error) {
	if _, exist := repo.sessionUserIdMap[id]; !exist {
		return []*domain.Session{}, nil
	}

	return slices.Collect[*domain.Session](maps.Values(repo.sessionUserIdMap[id])), nil
}
