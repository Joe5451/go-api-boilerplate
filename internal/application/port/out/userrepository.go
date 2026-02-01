package out

import (
	"context"
	"go-api-boilerplate/internal/domain"
)

type UserRepository interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user domain.User) error

	// GetUserByID gets a user by id
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)

	// GetUserByEmail gets a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
