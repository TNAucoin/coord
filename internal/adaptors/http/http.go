package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnaucoin/coord/config"
	"github.com/tnaucoin/coord/internal/ports"
)

type Adapter struct {
	api  ports.APIPort
	port string
}

func NewAdapter(api ports.APIPort) *Adapter {
	port := fmt.Sprintf(":%d", config.GetApplicationPort())
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	router := gin.Default()
	router.POST("/job", a.SubmitJob)
	router.Run(a.port)

}

func (a Adapter) SubmitJob(c *gin.Context) {
	c.JSON(200, nil)
}
