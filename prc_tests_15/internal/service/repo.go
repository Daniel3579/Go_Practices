package service

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidID    = errors.New("invalid user id")
	ErrEmptyName    = errors.New("name cannot be empty")
)

type User struct {
	ID        int64
	Email     string
	Name      string
	CreatedAt time.Time
}

type UserRepository interface {
	GetUser(id int64) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(u User) (int64, error)
	DeleteUser(id int64) error
}
