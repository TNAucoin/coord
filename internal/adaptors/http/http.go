package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	r.HandleFunc()

	srv := &http.Server{
		Addr:    a.port,
		Handler: r,
	}

	srv.ListenAndServe()
}
