package repository

import (
	"github.com/tnaucoin/coord/internal/core/dto"
	"github.com/tnaucoin/coord/internal/core/port/repository"
)

type jobRepository struct {
	db repository.Database
}

func NewJobRepository(db repository.Database) repository.JobRepository {
	return &jobRepository{
		db: db,
	}
}

func (j jobRepository) Insert(job dto.JobTask) error {
	//TODO: handle inserting the job record into the db
	return nil
}
