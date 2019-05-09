package storage

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	uuid "github.com/satori/go.uuid"
)

var (
	once sync.Once
	es   *elasticsearch.Client
)

// InitES initializes the Elasticsearch (Cluster) connection
func InitES(cfg elasticsearch.Config) {
	var err error
	es, err = elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatal("Could not connect to Elasticsearch. Check connection")
	}

	// Get elasticsearch cluster info
	res, err := es.Info()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Elasticsearch failure: %s", err)
	}

	log.Printf("Connected to elasticsearch with client version %s", elasticsearch.Version)
}

// AddDocument @TODO
func AddDocument(d string) error {
	c := make(chan error)

	// @TODO, should not be needed atm but good for multithreading later
	go func() {
		req := esapi.IndexRequest{
			Index:      "esfeed",              // Where to store it. Maybe make it dependend on the user instead of storing all data on a common index @TODO
			DocumentID: uuid.NewV4().String(), // Not sure if a uuid is needed here/how to generate it on the storage side.
			Body:       strings.NewReader(d),  // JSON Body
			Refresh:    "true",                // Refresh the index, maybe call this periodically instead
		}

		// Insert into elasticsearch
		res, err := req.Do(context.Background(), es)

		if err != nil {
			c <- err
			return
		}

		if res.IsError() {
			c <- errors.New(res.String())
			return
		}

		log.Print("Inserted document into elasticsearch")
		c <- nil
	}()

	return <-c
}
