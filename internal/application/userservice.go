package application

import (
	"context"
	"go-api-boilerplate/internal/application/port/in"
	"go-api-boilerplate/internal/application/port/out"
	"go-api-boilerplate/internal/domain"
)

type UserService struct {
	userRepo out.UserRepository
}

var _ in.UserUseCase = &UserService{}

func NewUserService(userRepo out.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, displayName, email, password string) error {
	user, err := domain.RegisterUser(displayName, email, password)
	if err != nil {
		return err
	}

	return s.userRepo.CreateUser(ctx, *user)
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := user.Authenticate(password); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
