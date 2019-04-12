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
	created    time.Time `json:"created"`
	Caption    string    `json:"caption"`
	Content    string    `json:"content"`
	Url        string    `json:"url"`
	IsFromFeed bool      `json:"isFromFeed"`
	FeedUrl    string    `json:"feedUrl"`
}

func (d Document) Dump() ([]byte, error) {
	return json.Marshal(d)
}
