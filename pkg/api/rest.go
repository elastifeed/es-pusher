package api

import (
	"encoding/json"
	"net/http"

	"github.com/elastifeed/es-pusher/pkg/document"
	"github.com/elastifeed/es-pusher/pkg/storage"
	"github.com/golang/glog"
)

// Restr is the REST API Interface which provides all endpoints
type Restr interface {
	AddDocuments(w http.ResponseWriter, r *http.Request)
}

// rests is the internal storage struct for the REST API. It contains the storage engine
type rests struct {
	storage storage.Storager
}

// New Creates a new REST API endpoint
func New(s storage.Storager) Restr {
	return rests{storage: s}
}

// AddDocuments adds 1..n documents to the elasticsearch database
func (rs rests) AddDocuments(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Index string              `json:"index"`
		Docs  []document.Document `json:"docs"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Has to match

	if decoder.Decode(&req) != nil {
		glog.Error("Error decoding Document from JSON Body")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("{\"status\": \"bad request\"}"))
		if err != nil {
			glog.Error("Response not fully transmitted")
		}
		return
	}

	err := rs.storage.AddDocuments(req.Index, req.Docs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("{\"status\": \"internal error\"}"))
		if err != nil {
			glog.Error("Response not fully transmitted")
		}
		return
	}

	_, err = w.Write([]byte("{\"status\": \"ok\"}"))
	if err != nil {
		glog.Error("Response not fully transmitted")
	}
}
