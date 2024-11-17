package repository

import (
	"context"
	"fmt"

	"github.com/gera9/user-service/pkg/models"
	"github.com/gera9/user-service/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type usersRepository struct {
	conn *pgxpool.Pool
}

func NewUsersRepository(conn *pgxpool.Pool) *usersRepository {
	return &usersRepository{conn: conn}
}

func (r *usersRepository) Register(ctx context.Context, user models.UserPayload) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := r.conn.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *usersRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at FROM users WHERE username = $1`
	err := r.conn.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *usersRepository) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`
	err := r.conn.QueryRow(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *usersRepository) UpdateById(ctx context.Context, id uuid.UUID, user models.UserPayload) error {
	query, pos, args := utils.UpdateQueryBuilder("users", user)
	query += fmt.Sprintf(" WHERE id = $%d", pos)
	args = append(args, id)

	_, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *usersRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
