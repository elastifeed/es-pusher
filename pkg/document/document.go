package document

import (
	"time"

	"encoding/json"
)

// Dumper interfaces specifies the dump function which converts the parsed struct
// into json again
type Dumper interface {
	Dump() string
}

// Document contains the schema of a data entry which is stored in Elasticsearch
type Document struct {
	Created    time.Time `json:"created"`
	Caption    string    `json:"caption"`
	Content    string    `json:"content"`
	URL        string    `json:"url"`
	IsFromFeed bool      `json:"isFromFeed"`
	FeedURL    string    `json:"feedUrl"`
}

// Dump Creates JSON Representation from a Elasticsearch Document
func (d Document) Dump() ([]byte, error) {
	return json.Marshal(d)
}

// Load loads a Document from a JSON Representation
func Load(data []byte) (Document, error) {
	d := Document{}
	err := json.Unmarshal(data, &d)

	return d, err
}
