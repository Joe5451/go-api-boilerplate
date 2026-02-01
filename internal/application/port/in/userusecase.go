package in

import (
	"context"
	"go-api-boilerplate/internal/domain"
)

type UserUseCase interface {
	// RegisterUser registers a new user
	RegisterUser(ctx context.Context, displayName, email, password string) error

	// AuthenticateUser authenticates a user
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)

	// GetUserProfile gets a user profile by id
	GetUserProfile(ctx context.Context, id int64) (*domain.User, error)
}
