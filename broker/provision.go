package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/types-cf"
	"github.com/go-chi/chi"
)

func Provision(r chi.Router) {
	// Status for provision request
	r.Get("/lastOperation", lastOperation)

	// Create service instance
	r.Put("/", create)

	// Delete service instance
	r.Delete("/", delete)
}

func lastOperation(w http.ResponseWriter, r *http.Request) {
	lastOp := cf.ServiceLastOperationResponse{
		State:       "succeeded",
		Description: "Service ready to use",
	}

	json, err := json.Marshal(lastOp)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this service instance:")
		//fmt.Printf("%# v\n", pretty.Formatter(lastOp))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func create(w http.ResponseWriter, r *http.Request) {
	// This app is also the service, so no need to create an actual service instance
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func delete(w http.ResponseWriter, r *http.Request) {
	// This app is also the service, so no need to delete an actual service instance
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
