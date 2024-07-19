package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnaucoin/coord/internal/core/common/router"
	"github.com/tnaucoin/coord/internal/core/model/request"
	"github.com/tnaucoin/coord/internal/core/port/service"
)

type JobController struct {
	gin        *gin.Engine
	jobService service.JobService
}

func NewJobController(
	gin *gin.Engine,
	jobService service.JobService,
) JobController {
	return JobController{
		gin:        gin,
		jobService: jobService,
	}
}

func (j JobController) InitRouter() {
	api := j.gin.Group("/api/v1")
	router.Post(api, "job", j.submitJob)
}

func (j JobController) submitJob(c *gin.Context) {
	log.Println("request hit")
	//TODO: handle the request
	_, err := j.parseRequest(c)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, "hello")
}

func (j JobController) parseRequest(ctx *gin.Context) (*request.SubmitJobRequest, error) {
	var req request.SubmitJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
