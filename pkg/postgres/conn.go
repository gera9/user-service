package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once sync.Once
	conn *pgxpool.Pool
)

func NewPostgresConn(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	if conn != nil {
		return conn, nil
	}

	var (
		dbpool *pgxpool.Pool
		err    error
	)
	once.Do(func() {
		dbpool, err = pgxpool.New(ctx, connString)
		if err != nil {
			return
		}

		err = dbpool.Ping(ctx)
		if err != nil {
			return
		}

		conn = dbpool
	})

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetConn() *pgxpool.Pool {
	return conn
}
