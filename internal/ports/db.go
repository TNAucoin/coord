package ports

import "github.com/tnaucoin/coord/internal/application/core/domain"

type APIPort interface {
	SubmitJob(job domain.Job) (domain.Job, error)
}

type DBPort interface {
	GetJob(id string) (domain.Job, error)
	SaveJob(*domain.Job) error
}
