package main

import (
	"log"

	"github.com/tnaucoin/coord/config"
	"github.com/tnaucoin/coord/internal/adaptors/db"
	"github.com/tnaucoin/coord/internal/adaptors/http"
	"github.com/tnaucoin/coord/internal/application/core/api"
)

func main() {
	dbAdaptor, err := db.NewAdaptor(config.GetDatabaseConnectionURL())
	if err != nil {
		log.Fatalf("failed to connect to the db: %v", err)
	}
	application := api.NewApplication(dbAdaptor)
	httpAdaptor := http.NewAdapter(application)
	httpAdaptor.Run()
}
