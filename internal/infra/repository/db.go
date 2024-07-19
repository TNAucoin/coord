package repository

import (
	"database/sql"

	"github.com/tnaucoin/coord/internal/core/port/repository"
)

type database struct {
	*sql.DB
}

func NewDB() (repository.Database, error) {
	//TODO: implement db connection here
	return nil, nil
}

func (d database) GetDB() *sql.DB {
	return d.DB
}
