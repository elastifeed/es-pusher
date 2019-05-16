package storage

import "github.com/elastifeed/es-pusher/pkg/document"

// Storage interface
type Storage interface {
	AddDocuments(index string, docs []document.Document) error
}
