package main

import (
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastifeed/es-pusher/pkg/api"
	"github.com/elastifeed/es-pusher/pkg/storage"
)

func main() {

	// Connect to a specified Elasticsearch instance
	s := storage.NewES(elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	})

	// Create new Rest Api Endpoint based on the previously connected elasticsearch storage engine
	rAPI := api.New(s)
	// Add HTTP Endpoint to /add
	http.HandleFunc("/add", rAPI.AddDocuments)

	// Run forever and exit on error
	log.Fatal(http.ListenAndServe(":8080", nil))
}
