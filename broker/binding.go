package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func Binding(r chi.Router) {
	// Bind service instance
	r.Put("/", bind)

	// Unbind service instance
	r.Delete("/", unbind)
}

func bind(w http.ResponseWriter, r *http.Request) {
	type ServiceBindingResonse struct {
		RouteServiceUrl string `json:"route_service_url"`
	}
	bind := ServiceBindingResonse{
		RouteServiceUrl: "https://" + r.Host,
	}

	json, err := json.Marshal(bind)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this service binding response")
		//fmt.Printf("%# v\n", pretty.Formatter(lastOp))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func unbind(w http.ResponseWriter, r *http.Request) {
	// No binding resources to be deleted
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
