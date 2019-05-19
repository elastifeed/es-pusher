package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastifeed/es-pusher/pkg/api"
	"github.com/elastifeed/es-pusher/pkg/storage"
)

// main entry point for es-pusher
//
// Possible Configuration options (via environment):
//    - ES_ADDRESSES (e.g. ["http://host0:9200", "http://host1:9200", ...])
//    - API_BIND (e.g. ":9000")
func main() {

	var addrs []string

	if json.Unmarshal([]byte(os.Getenv("ES_ADDRESSES")), &addrs) != nil {
		log.Fatal("Configuration error, check ES_ADDRESSES")
	}

	log.Println(addrs)

	// Connect to a specified Elasticsearch instance
	s := storage.NewES(elasticsearch.Config{
		Addresses: addrs,
	})

	// Create new Rest Api Endpoint based on the previously connected elasticsearch storage engine
	rAPI := api.New(s)
	// Add HTTP Endpoint to /add
	http.HandleFunc("/add", rAPI.AddDocuments)

	// Run forever and exit on error
	log.Fatal(http.ListenAndServe(os.Getenv("API_BIND"), nil))
}
