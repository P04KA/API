package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConnect(connStr string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrap(err, "pgx conn")
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "conn ping")
	}
	return conn, nil
}
