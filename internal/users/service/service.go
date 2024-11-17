package service

import (
	"context"
	"errors"

	"github.com/gera9/user-service/config"
	"github.com/gera9/user-service/internal/users"
	"github.com/gera9/user-service/pkg/models"
	"github.com/gera9/user-service/pkg/utils"
	"github.com/google/uuid"
)

type usersService struct {
	repo users.UserRepository
}

func NewUsersService(repo users.UserRepository) *usersService {
	return &usersService{
		repo: repo,
	}
}

func (u *usersService) DeleteById(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteById(ctx, id)
}

func (u *usersService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return u.repo.GetByUsername(ctx, username)
}

func (u *usersService) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return u.repo.GetById(ctx, id)
}

func (u *usersService) LoginByUsername(ctx context.Context, user models.UserPayload) (string, error) {
	if user.Username == "" || user.Password == "" {
		return "", errors.New("username and password are required")
	}

	dbUser, err := u.repo.GetByUsername(ctx, user.Username)
	if err != nil {
		return "", err
	}

	if !dbUser.ComparePassword(user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := utils.CreateToken(config.Config{}, dbUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *usersService) Register(ctx context.Context, user models.UserPayload) (uuid.UUID, error) {
	err := user.Validate()
	if err != nil {
		return uuid.Nil, err
	}

	err = user.HashPassword()
	if err != nil {
		return uuid.Nil, err
	}

	return u.repo.Register(ctx, user)
}

func (u *usersService) UpdateById(ctx context.Context, id uuid.UUID, user models.UserPayload) error {
	return u.repo.UpdateById(ctx, id, user)
}
