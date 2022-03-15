package main

import (
	"learn-typesense/internal/core/service"
	"learn-typesense/internal/core/ts"
	"learn-typesense/internal/driven/storage"
	"learn-typesense/internal/driver/rest"
	"log"
	"net/http"

	"github.com/typesense/typesense-go/typesense"
)

func main() {
	// initialize typesense client
	tsClient := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("Hu52dwsas2AdxdE"),
	)

	// initialize typesense storage
	tsStorage, err := ts.NewStorage(
		ts.StorageConfig{
			TSClient: tsClient,
		},
	)
	if err != nil {
		log.Fatalf("unable to initialize typesense storage due: %w", err)
	}

	// initialize book schema in typesense storage
	// defaultSort := "year"
	collectionName := "book"
	// fields := []api.Field{
	// 	{
	// 		Name: "id",
	// 		Type: "string",
	// 	},
	// 	{
	// 		Name: "title",
	// 		Type: "string",
	// 	},
	// 	{
	// 		Name: "author",
	// 		Type: "string",
	// 	},
	// 	{
	// 		Name: "year",
	// 		Type: "int32",
	// 	},
	// 	{
	// 		Name: "summary",
	// 		Type: "string",
	// 	},
	// }
	// err = tsStorage.CreateSchema(collectionName, fields, defaultSort)
	// if err != nil {
	// 	log.Fatalf("unable to initialize book schema due: %w", err)
	// }

	// initailize book storage
	bookStorage, err := storage.NewStorage(
		storage.StorageConfig{
			TSStorage:        tsStorage,
			TSCollectionName: collectionName,
		},
	)
	if err != nil {
		log.Fatalf("unable to initialize book storage due: %w", err)
	}

	// initialize service
	service, err := service.NewService(
		service.ServiceConfig{
			BookStorage: bookStorage,
		},
	)
	if err != nil {
		log.Fatalf("unable to initialize service due: %w", err)
	}

	// initialize rest api
	api := rest.NewAPI(rest.APIConfig{
		BookService: service,
	})

	// initialize server
	addr := "9000"
	server := &http.Server{
		Addr:    ":" + addr,
		Handler: api.Handler(),
	}

	// run server
	log.Printf("server is listening on %v...", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start server due: %v", err)
	}
}
