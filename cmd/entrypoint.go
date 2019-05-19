package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastifeed/es-pusher/pkg/api"
	"github.com/elastifeed/es-pusher/pkg/storage"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main entry point for es-pusher
//
// Possible Configuration options (via environment):
//    - ES_ADDRESSES (e.g. ["http://host0:9200", "http://host1:9200", ...])
//    - API_BIND (e.g. ":9000")
func main() {
	flag.Parse()

	var addrs []string

	if json.Unmarshal([]byte(os.Getenv("ES_ADDRESSES")), &addrs) != nil {
		glog.Fatal("Configuration error, check ES_ADDRESSES")
	}

	// Connect to a specified Elasticsearch instance
	s, err := storage.NewES(elasticsearch.Config{
		Addresses: addrs,
	})

	if err != nil {
		glog.Fatal("Elasticsearch Initialization failed", err)
	}

	// Create new Rest Api Endpoint based on the previously connected elasticsearch storage engine
	rAPI := api.New(s)
	// Add API Endpoint to /add
	http.HandleFunc("/add", rAPI.AddDocuments)

	// Add Monitoring endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Run forever and exit on error
	glog.Fatal(http.ListenAndServe(os.Getenv("API_BIND"), nil))
}
