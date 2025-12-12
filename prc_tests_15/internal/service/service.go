package service

import (
	"strings"
	"time"
)

type UserService struct {
	repo UserRepository
}

func New(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindIDByEmail(email string) (int64, error) {
	if email == "" {
		return 0, ErrInvalidEmail
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *UserService) FindByEmail(email string) (User, error) {
	if email == "" {
		return User{}, ErrInvalidEmail
	}

	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByID(id int64) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidID
	}

	return s.repo.GetUser(id)
}

func (s *UserService) CreateUser(email, name string) (int64, error) {

	if email == "" {
		return 0, ErrInvalidEmail
	}

	atIndex := strings.Index(email, "@")
	dotIndex := strings.LastIndex(email, ".")

	if atIndex <= 0 || dotIndex <= atIndex+1 || dotIndex >= len(email)-1 {
		return 0, ErrInvalidEmail
	}

	if strings.TrimSpace(name) == "" {
		return 0, ErrEmptyName
	}

	user := User{
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateUser(user)
}
