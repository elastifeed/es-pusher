package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	uuid "github.com/satori/go.uuid"
)

func main() {
	var wg sync.WaitGroup

	// Connection URL for Elasticsearch, localhost for development. @TODO
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	}

	// Connect
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Could not instantiate elasticsearch client: %s", err)
	}

	// Query info of elasticsearch
	res, err := es.Info()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Elasticsearch failure: %s", err)
	}

	log.Printf("Connected to server version %s", elasticsearch.Version)

	// Iterate over some test data, nothing too fancy @TODO
	for i, title := range []string{"Test 1", "Test 2", "Test 3"} {
		// Use a Waitgroup to ensure everything was sent to Elasticsearch before exiting the program
		wg.Add(1)
		go func(i int, title string) {
			// Gets executed when the function exits
			defer wg.Done()

			// UUID, might be better to use an integer though I am not sure how to lookup the last
			// index or generate it at the storage side @TODO
			id := uuid.NewV4().String()

			fmt.Println(id)

			// Direct object - template is definitely needed. Just a dummy for now! @TODO
			req := esapi.IndexRequest{
				Index:      "test",                                           // Where to store it. Maybe make it dependend on the user instead of storing all data on a common index @TODO
				DocumentID: id,                                               // Not sure if a uuid is needed here/how to generate it on the storage side.
				Body:       strings.NewReader(`{"title" : "` + title + `"}`), // JSON Body
				Refresh:    "true",                                           // Refresh the index, maybe call this periodically instead
			}

			// Perform request with the previous established connection
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close() // Be sure to close the response. defer is awesome!

			if res.IsError() {
				log.Printf("[%s] Could not index document with ID=%d", res.Status(), id)
			}
		}(i, title)

	}
	wg.Wait()
}
