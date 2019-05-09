package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastifeed/es-pusher/pkg/document"
	"github.com/elastifeed/es-pusher/pkg/storage"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var d document.Document

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Has to match

	if decoder.Decode(&d) != nil {
		log.Printf("Error decoding Document from JSON Body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"status\": \"bad request\"}"))
		return
	}

	res, err := d.Dump()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"status\": \"internal error\"}"))
		return
	}

	// Insert into DB
	err = storage.AddDocument(string(res))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"status\": \"internal error\"}"))
		return
	}

	w.Write([]byte("{\"status\": \"ok\"}"))
}

func main() {

	storage.InitES(elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	})

	http.HandleFunc("/add", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
