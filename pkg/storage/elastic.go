package storage

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastifeed/es-pusher/pkg/document"
	uuid "github.com/satori/go.uuid"
)

// esdriver
type esdriver struct {
	es *elasticsearch.Client
}

// NewES establishes a new Elasticsearch connection
func NewES(cfg elasticsearch.Config) Storage {
	var e esdriver
	var err error
	var r map[string]interface{}

	e.es, err = elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatal("Could not connect to Elasticsearch. Check connection")
	}

	// Get elasticsearch cluster info
	res, err := e.es.Info()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Elasticsearch failure: %s", err)
	}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the elasticsearch response body: %s", err)
	}

	log.Printf("Connected to elasticsearch %s", r["version"].(map[string]interface{})["number"])

	return e
}

// AddDocuments @TODO
func (e esdriver) AddDocuments(index string, docs []document.Document) error {
	var wg sync.WaitGroup

	// @TODO, should not be needed atm but good for multithreading later
	for _, d := range docs {
		dString, _ := d.Dump()
		go func(toAdd string) {
			wg.Add(1)
			defer wg.Done()

			req := esapi.IndexRequest{
				Index:      index,                    // Where to store it. Maybe make it dependend on the user instead of storing all data on a common index @TODO
				DocumentID: uuid.NewV4().String(),    // Not sure if a uuid is needed here/how to generate it on the storage side.
				Body:       strings.NewReader(toAdd), // JSON Body
				Refresh:    "true",                   // Refresh the index, maybe call this periodically instead
			}

			// Insert into elasticsearch
			res, err := req.Do(context.Background(), e.es)

			if err != nil {
				return
			}

			if res.IsError() {
				return
			}

			log.Print("Inserted document into elasticsearch")
		}(string(dString))
	}

	wg.Wait()

	return nil
}
