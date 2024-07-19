package repository

import "github.com/tnaucoin/coord/internal/core/dto"

type JobRepository interface {
	Insert(job dto.JobTask) error
}
