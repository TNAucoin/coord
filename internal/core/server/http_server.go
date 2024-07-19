package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnaucoin/coord/internal/core/config"
)

const defaultHost = "0.0.0.0"

type HttpServer interface {
	Start()
	Stop()
}

type httpServer struct {
	Port   uint
	server *http.Server
}

func NewHttpServer(router *gin.Engine, config config.HttpServerConfig) HttpServer {
	return &httpServer{
		Port: config.Port,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", defaultHost, config.Port),
			Handler: router,
		},
	}
}

func (h httpServer) Start() {
	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}()
	log.Printf("Started Job Service on port %d", h.Port)
}
func (h httpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err.Error())
	}
}
