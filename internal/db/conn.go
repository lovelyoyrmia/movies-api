package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lovelyoyrmia/movies-api/pkg/config"
)

type Database struct {
	DB *pgxpool.Pool
}

func NewDatabase(ctx context.Context, conf config.Config) (*Database, error) {
	sqlDriver, err := pgxpool.New(ctx, conf.DBUrl)
	if err != nil {
		return nil, err
	}
	return &Database{
		DB: sqlDriver,
	}, nil
}
