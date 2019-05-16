package storage

import "github.com/elastifeed/es-pusher/pkg/document"

// Stringer interface for documents
type Stringer interface {
	String() string
}

// Storage interface
type Storage interface {
	AddDocuments(index string, docs []document.Document) error
}
