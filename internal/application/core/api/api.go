package api

import (
	"fmt"

	"github.com/tnaucoin/coord/internal/application/core/domain"
	"github.com/tnaucoin/coord/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) SubmitJob(job *domain.Job) (*domain.Job, error) {
	err := a.db.SaveJob(job)
	if err != nil {
		return &domain.Job{}, fmt.Errorf("failed to save job: %v", err)
	}
	return job, nil
}
