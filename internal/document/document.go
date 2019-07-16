package document

import (
	"time"

	"encoding/json"
)

// Document contains the schema of a data entry which is stored in Elasticsearch
type Document struct {
	Created         time.Time `json:"created"`
	Author          string    `json:"author"`
	Title           string    `json:"title"`
	RawContent      string    `json:"raw_content"`
	MarkdownContent string    `json:"markdown_content"`
	PdfURL          string    `json:"pdf"`
	ScreenshotURL   string    `json:"screenshot"`
	ThumbnailURL    string    `json:"thumbnail"`
	URL             string    `json:"url"`
	Categories      []string  `json:"categories"`
	IsFromFeed      bool      `json:"from_feed"`
	FeedURL         string    `json:"feed_url"`
	Starred         bool      `json:"starred"`
	ReadLater       bool      `json:"read_later"`
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
