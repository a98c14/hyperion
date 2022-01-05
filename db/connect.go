package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func InitDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(ConnectionString)
	if err != nil {
		panic("Heck1!")
	}
	config.MaxConns = 1000
	config.MaxConnLifetime = 180
	config.MaxConnIdleTime = 180

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		panic("Heck2!")
	}
	Pool = pool
	return pool, err
}

func GetConnectionPool(ctx context.Context) (*pgxpool.Pool, error) {
	return Pool, nil
}
