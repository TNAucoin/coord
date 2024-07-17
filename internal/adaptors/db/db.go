package db

import (
	"fmt"

	"github.com/tnaucoin/coord/internal/application/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Status      string
	JobSteps    []JobStep `gorm:"foreignKey:JobID"`
	CurrentStep int
}

type JobStep struct {
	gorm.Model
	// Args map[string]interface{}
	Type  string
	JobID uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Adapter struct {
	db *gorm.DB
}

func NewAdaptor(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(postgres.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	err := db.AutoMigrate(&Job{}, &JobStep{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) GetJob(id string) (domain.Job, error) {
	var jobEntity domain.Job
	res := a.db.First(&jobEntity, id)
	var jobSteps []domain.JobStep
	for _, jobEntity := range jobEntity.JobSteps {
		jobSteps = append(jobSteps, domain.JobStep{
			Type: jobEntity.Type,
		})
	}
	job := domain.Job{
		ID:          int64(jobEntity.ID),
		Status:      jobEntity.Status,
		JobSteps:    jobSteps,
		CurrentStep: jobEntity.CurrentStep,
	}
	return job, res.Error
}

func (a Adapter) SaveJob(job *domain.Job) error {
	var jobSteps []JobStep
	for _, jobStep := range job.JobSteps {
		jobSteps = append(jobSteps, JobStep{
			Type: jobStep.Type,
		})
	}
	jobModel := Job{
		Status:      job.Status,
		JobSteps:    jobSteps,
		CurrentStep: job.CurrentStep,
	}
	res := a.db.Create(&jobModel)
	if res.Error != nil {
		job.ID = int64(jobModel.ID)
	}
	return res.Error
}
