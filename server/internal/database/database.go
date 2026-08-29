package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect(ctx context.Context, databaseURL string) error {
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return err
	}
	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		pool = nil
		return err
	}
	return nil
}

func Pool() *pgxpool.Pool {
	if pool == nil {
		panic("database pool used before initialization")
	}
	return pool
}

func Close() {
	if pool != nil {
		pool.Close()
		pool = nil
	}
}
