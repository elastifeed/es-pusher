package main

import (
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastifeed/es-pusher/pkg/api"
	"github.com/elastifeed/es-pusher/pkg/storage"
)

func main() {

	s := storage.NewES(elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	})

	if s == nil {
		return
	}

	r := api.New(s)

	http.HandleFunc("/add", r.AddDocuments)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
