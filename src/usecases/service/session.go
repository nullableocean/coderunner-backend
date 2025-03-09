package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"
	"nullableocean-postupashki/src/usecases"
)

var (
	sessionIdLength = 32
)

type SessionManager struct {
	sRepo repository.Session
}

func NewSessionService(sRepo repository.Session) usecases.Session {
	return &SessionManager{
		sRepo: sRepo,
	}
}

func (s *SessionManager) Start(user *domain.User) (*domain.Session, error) {
	if user.Id == 0 {
		return nil, errors.New("need user id")
	}

	sessionId, err := s.createSessionId()
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		UserId:    user.Id,
		SessionId: sessionId,
	}

	s.sRepo.Save(session)

	return session, nil
}

func (s *SessionManager) Get(sid string) (*domain.Session, error) {
	return s.sRepo.Get(sid)
}

func (s *SessionManager) Destroy(sid string) error {
	return s.sRepo.Delete(sid)
}

func (s *SessionManager) DestroyAllByUser(userId int64) error {
	sessions, err := s.sRepo.GetByUserId(userId)
	if err != nil {
		return err
	}

	for _, ses := range sessions {
		if err := s.sRepo.Delete(ses.SessionId); err != nil {
			return err
		}
	}

	return nil
}

func (s *SessionManager) createSessionId() (string, error) {
	b := make([]byte, sessionIdLength)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
