package service

import (
	"net/http"

	"github.com/tnaucoin/coord/internal/core/model/request"
	"github.com/tnaucoin/coord/internal/core/model/response"
	"github.com/tnaucoin/coord/internal/core/port/repository"
	"github.com/tnaucoin/coord/internal/core/port/service"
)

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) service.JobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}

func (j jobService) SubmitJob(request *request.SubmitJobRequest) *response.Response {
	//TODO: validate / handle job request
	//TODO: use jobRepository to insert the job
	//TODO: replace this with a real job resp
	return j.createSuccessResponse(response.SubmitJobResponse{})
}

func (j jobService) createSuccessResponse(data response.SubmitJobResponse) *response.Response {
	return &response.Response{
		Data:    data,
		Status:  true,
		Code:    http.StatusCreated,
		Message: "job submitted successfully",
	}
}

func (j jobService) createFailedResponse(code int, message string) *response.Response {
	return &response.Response{
		Data:    nil,
		Status:  false,
		Code:    code,
		Message: message,
	}
}
