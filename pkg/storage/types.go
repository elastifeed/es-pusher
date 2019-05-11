package storage

// Stringer interface for documents
type Stringer interface {
	String() string
}

// Storage interface
type Storage interface {
	AddDocuments(index string, docs []Stringer) error
}
