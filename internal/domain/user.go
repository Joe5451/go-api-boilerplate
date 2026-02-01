package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDisplayNameRequired = errors.New("display name is required")
	ErrEmailRequired       = errors.New("email is required")
	ErrDuplicatedEmail     = errors.New("email already exists")
	ErrPasswordRequired    = errors.New("password is required")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)

type User struct {
	ID           int64     `json:"id"`
	DisplayName  string    `json:"display_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func RegisterUser(displayName, email, password string) (*User, error) {
	if displayName == "" {
		return nil, ErrDisplayNameRequired
	}

	if email == "" {
		return nil, ErrEmailRequired
	}

	if password == "" {
		return nil, ErrPasswordRequired
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func (u *User) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
