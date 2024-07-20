package repository

import (
	"github.com/tnaucoin/coord/coord"
	"github.com/tnaucoin/coord/internal/core/dto"
	"github.com/tnaucoin/coord/internal/core/port/repository"
)

type jobRepository struct {
	jobDB *coord.Queries
}

func NewJobRepository(db repository.Database) repository.JobRepository {
	jobDB := coord.New(db.GetDB())
	return &jobRepository{
		jobDB: jobDB,
	}
}

func (j jobRepository) Insert(job dto.JobTask) error {
	//TODO: handle inserting the job record into the db
	j.jobDB.InsertJob()
	return nil
}
