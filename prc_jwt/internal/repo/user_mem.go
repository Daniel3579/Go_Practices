package repo

import (
	"errors"

	"example.com/prc_jwt/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserMem struct{ users map[string]core.User } // key = email

func NewUserMem() *UserMem {
	// заранее захэшированные пароли (пример: "secret123")
	hash := func(s string) []byte { h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost); return h }
	return &UserMem{users: map[string]core.User{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
	}}
}

var ErrNotFound = errors.New("user not found")
var ErrBadCreds = errors.New("bad credentials")

func (r *UserMem) ByEmail(email string) (core.User, error) {
	u, ok := r.users[email]
	if !ok {
		return core.User{}, ErrNotFound
	}
	return u, nil
}

func (r *UserMem) CheckPassword(email, pass string) (core.User, error) {
	u, err := r.ByEmail(email)
	if err != nil {
		return core.User{}, ErrNotFound
	}
	if bcrypt.CompareHashAndPassword(u.Hash, []byte(pass)) != nil {
		return core.User{}, ErrBadCreds
	}
	return u, nil
}
