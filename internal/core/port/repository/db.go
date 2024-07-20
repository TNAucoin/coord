package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Database interface {
	GetDB() *pgx.Conn
	Close(ctx context.Context) error
}
