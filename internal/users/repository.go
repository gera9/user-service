package users

import (
	"context"

	"github.com/gera9/user-service/pkg/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Register(ctx context.Context, user models.UserPayload) (uuid.UUID, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateById(ctx context.Context, id uuid.UUID, user models.UserPayload) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
