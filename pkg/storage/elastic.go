package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastifeed/es-pusher/pkg/document"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"k8s.io/klog"
)

// esdriver
type esdriver struct {
	es *elasticsearch.Client
}

var (
	addDocumentRequestedCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "espusher_storage_elasticsearch_added_document_total",
		Help: "Number of Documents added to elasticsearch",
	})
	addedDocumentCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "espusher_storage_elasticsearch_add_document_requested_total",
		Help: "Number of Documents added to elasticsearch",
	})
)

// NewES establishes a new Elasticsearch connection
func NewES(cfg elasticsearch.Config) (Storager, error) {
	var e esdriver
	var err error
	var r map[string]interface{}

	e.es, err = elasticsearch.NewClient(cfg)

	if err != nil {
		klog.Fatal("Could not connect to Elasticsearch. Check connection")
	}

	// Get elasticsearch cluster info
	res, err := e.es.Info()

	if err != nil {
		klog.Errorf("Error getting response: %s", err)
		return nil, err
	}

	if res.IsError() {
		klog.Errorf("Elasticsearch failure: %s", err)
		return nil, err
	}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		klog.Errorf("Error parsing the elasticsearch response body: %s", err)
		return nil, err
	}

	klog.Infof("Connected to elasticsearch %s", r["version"].(map[string]interface{})["number"])

	return e, nil
}

// AddDocuments adds 1..n documents to elasticsearch.
func (e esdriver) AddDocuments(indexes []string, docs []document.Document) error {
	var wg sync.WaitGroup

	// Update counter
	addDocumentRequestedCount.Add(float64(len(docs)))

	// Add all documents
	for _, d := range docs {

		wg.Add(1)
		go func(d document.Document) {
			// Generate hashed document index to avoid duplicates
			idHasher := sha256.New()
			io.WriteString(idHasher, d.URL+d.Caption+d.Content)
			dString, _ := d.Dump()
			defer wg.Done()

			for _, index := range indexes {
				wg.Add(1)
				go func(index string) {
					defer wg.Done()
					req := esapi.IndexRequest{
						Index:      index,
						DocumentID: hex.EncodeToString(idHasher.Sum(nil)),
						Body:       strings.NewReader(string(dString)),
						Refresh:    "true",
					}

					// Insert into elasticsearch
					res, err := req.Do(context.Background(), e.es)

					if err != nil || res.IsError() {
						return
					}

					addedDocumentCount.Inc()
				}(index)
			}

		}(d)
	}

	wg.Wait()

	klog.Infof("Inserted %d documents into elasticsearch into \"%s\"", len(docs), indexes)
	return nil
}
