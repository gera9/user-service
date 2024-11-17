package utils

import "context"

type CRUD[T any] interface {
	GetAll(ctx context.Context, search string, limit, offset int) ([]T, error)
	GetByID(ctx context.Context, id int) (*T, error)
	Create(ctx context.Context, t T) error
	Update(ctx context.Context, t T) error
	Delete(ctx context.Context, id int) error
}
