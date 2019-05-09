package document

import (
	"time"

	"encoding/json"
)

type Dumper interface {
	Dump() string
}

// @TODO
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
