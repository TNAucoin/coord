package service

import (
	"github.com/tnaucoin/coord/internal/core/model/request"
	"github.com/tnaucoin/coord/internal/core/model/response"
)

type JobService interface {
	SubmitJob(request *request.SubmitJobRequest) *response.Response
}
