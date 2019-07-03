package api

import (
	"encoding/json"
	"net/http"

	"github.com/elastifeed/es-pusher/pkg/document"
	"github.com/elastifeed/es-pusher/pkg/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"k8s.io/klog"
)

// Restr is the REST API Interface which provides all endpoints
type Restr interface {
	AddDocuments(w http.ResponseWriter, r *http.Request)
}

// rests is the internal storage struct for the REST API. It contains the storage engine
type rests struct {
	storage storage.Storager
}

var (
	restCallsCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "espusher_rest_total",
		Help: "Number of REST-API calls",
	})
	restCallsMalformed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "espusher_rest_calls_malformed",
		Help: "Number of malformed (json) REST-API calls",
	})
	restCallsSuccessful = promauto.NewCounter(prometheus.CounterOpts{
		Name: "espusher_rest_calls_successful",
		Help: "Number of successful REST-API calls",
	})
)

// New Creates a new REST API endpoint
func New(s storage.Storager) Restr {
	return rests{storage: s}
}

// AddDocuments adds 1..n documents to the elasticsearch database
func (rs rests) AddDocuments(w http.ResponseWriter, r *http.Request) {
	restCallsCount.Inc()

	// Always JSON response here
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Index string              `json:"index"`
		Docs  []document.Document `json:"docs"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Has to match

	if decoder.Decode(&req) != nil {
		restCallsMalformed.Inc()
		klog.Error("Error decoding Document from JSON Body")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("{\"status\": \"bad request\"}"))
		if err != nil {
			klog.Error("Response not fully transmitted")
		}
		return
	}

	err := rs.storage.AddDocuments(req.Index, req.Docs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("{\"status\": \"internal error\"}"))
		if err != nil {
			klog.Error("Response not fully transmitted")
		}
		return
	}

	_, err = w.Write([]byte("{\"status\": \"ok\"}"))
	if err != nil {
		klog.Error("Response not fully transmitted")
	}

	restCallsSuccessful.Inc()
}
