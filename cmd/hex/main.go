package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tnaucoin/coord/internal/controller/http"
	"github.com/tnaucoin/coord/internal/core/config"
	"github.com/tnaucoin/coord/internal/core/server"
	"github.com/tnaucoin/coord/internal/core/service"
	"github.com/tnaucoin/coord/internal/infra/repository"
)

func main() {
	instance := gin.New()
	instance.Use(gin.Recovery())

	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("failed to init db connection: %s", err.Error())
	}

	jobRepo := repository.NewJobRepository(db)
	jobService := service.NewJobService(jobRepo)
	jobController := controller.NewJobController(instance, jobService)
	jobController.InitRouter()

	httpServer := server.NewHttpServer(instance, config.HttpServerConfig{Port: 8000})

	httpServer.Start()
	defer httpServer.Stop()

	//OS signal
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	<-c
	log.Println("graceful shutdown.")

}
