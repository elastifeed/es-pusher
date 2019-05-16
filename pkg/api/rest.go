package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elastifeed/es-pusher/pkg/document"
	"github.com/elastifeed/es-pusher/pkg/storage"
)

// Rest @TODO
type Rest interface {
	AddDocuments(w http.ResponseWriter, r *http.Request)
}

type rests struct {
	storage storage.Storage
}

// New @TODO
func New(s storage.Storage) Rest {
	return rests{storage: s}
}

// AddDocuments @TODO
func (rs rests) AddDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []document.Document

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Has to match

	if decoder.Decode(&docs) != nil {
		log.Printf("Error decoding Document from JSON Body")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("{\"status\": \"bad request\"}"))
		if err != nil {
			log.Print("Response not fully transmitted")
		}
		return
	}

	err := rs.storage.AddDocuments("testindex", docs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("{\"status\": \"internal error\"}"))
		if err != nil {
			log.Print("Response not fully transmitted")
		}
		return
	}

	_, err = w.Write([]byte("{\"status\": \"ok\"}"))
	if err != nil {
		log.Print("Response not fully transmitted")
	}
}
