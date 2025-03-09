package service

import (
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/pkg/hasher"
	"nullableocean-postupashki/src/repository"
	"nullableocean-postupashki/src/usecases"
)

type UserService struct {
	userRepo        repository.UserRepository
	sessionsService usecases.Session
	passHasher      hasher.PasswordHasher
}

func NewUserService(
	userRepo repository.UserRepository,
	sessionsService usecases.Session,
	passHasher hasher.PasswordHasher,
) usecases.User {
	return &UserService{
		userRepo:        userRepo,
		sessionsService: sessionsService,
		passHasher:      passHasher,
	}
}

func (s *UserService) Create(user domain.User) (*domain.User, error) {
	hashedPass, err := s.passHasher.Hash([]byte(user.Password))
	if err != nil {
		return nil, err
	}

	user.HashedPassword = string(hashedPass)
	user.Password = ""

	return s.userRepo.Create(&user)
}

func (s *UserService) SignIn(user *domain.User) (*domain.Session, error) {
	foundUser, err := s.userRepo.GetByLogin(user.Login)
	if err != nil {
		return nil, err
	}

	err = s.VerifyUser(user, foundUser)
	if err != nil {
		return nil, err
	}

	return s.CreateSession(foundUser)
}

func (s *UserService) CreateSession(user *domain.User) (*domain.Session, error) {
	return s.sessionsService.Start(user)
}

func (s *UserService) VerifyUser(verifiable, registered *domain.User) error {
	passVerifed, err := s.passHasher.Verify([]byte(verifiable.Password), []byte(registered.HashedPassword))
	if !passVerifed {
		return usecases.ErrUserPassNotVerified
	}

	return err
}
