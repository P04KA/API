package storage

import "github.com/jackc/pgx/v5/pgxpool"

func GetConnect() (*pgxpool.Pool, error)
