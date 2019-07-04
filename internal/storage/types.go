package storage

import "github.com/elastifeed/es-pusher/internal/document"

// Storager is contains all functions provided by the storage driver engine
type Storager interface {
	AddDocuments(index []string, docs []document.Document) error
}
