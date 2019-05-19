package storage

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastifeed/es-pusher/pkg/document"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

// esdriver
type esdriver struct {
	es *elasticsearch.Client
}

// NewES establishes a new Elasticsearch connection
func NewES(cfg elasticsearch.Config) (Storager, error) {
	var e esdriver
	var err error
	var r map[string]interface{}

	e.es, err = elasticsearch.NewClient(cfg)

	if err != nil {
		glog.Fatal("Could not connect to Elasticsearch. Check connection")
	}

	// Get elasticsearch cluster info
	res, err := e.es.Info()

	if err != nil {
		glog.Errorf("Error getting response: %s", err)
		return nil, err
	}

	if res.IsError() {
		glog.Errorf("Elasticsearch failure: %s", err)
		return nil, err
	}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		glog.Errorf("Error parsing the elasticsearch response body: %s", err)
		return nil, err
	}

	glog.Infof("Connected to elasticsearch %s", r["version"].(map[string]interface{})["number"])

	return e, nil
}

// AddDocuments adds 1..n documents to elasticsearch.
func (e esdriver) AddDocuments(index string, docs []document.Document) error {
	var wg sync.WaitGroup

	// Add all documents
	for _, d := range docs {
		dString, _ := d.Dump()
		wg.Add(1)
		go func(toAdd string) {
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
		}(string(dString))
	}

	wg.Wait()

	glog.Infof("Inserted %d documents into elasticsearch into \"%s\"", len(docs), index)
	return nil
}
