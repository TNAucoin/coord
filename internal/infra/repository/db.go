package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tnaucoin/coord/config"
	"github.com/tnaucoin/coord/internal/core/port/repository"
)

type database struct {
	DB *pgx.Conn
}

func NewDB(ctx context.Context) (repository.Database, error) {
	conn, err := pgx.Connect(ctx, config.GetDatabaseConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %s", err.Error())
	}
	db := &database{
		DB: conn,
	}
	return db, nil
}

func (d database) GetDB() *pgx.Conn {
	return d.DB
}

func (d database) Close(ctx context.Context) error {
	return d.DB.Close(ctx)
}
