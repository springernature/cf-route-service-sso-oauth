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
	// Check which oauth provider is requested to bind
	provider := "google" // default
	if p := bindingParameters(r)["provider"]; p != "" {
		provider = p
	}
	// Set the route service url
	bind := ServiceBindingResonse{
		RouteServiceUrl: "https://" + r.Host + "/" + provider,
	}
	json, err := json.Marshal(bind)
	if err != nil {
		fmt.Println("Failed to marshal the service binding response")
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

func bindingParameters(r *http.Request) map[string]string {
	type ServiceBindingRequest struct {
		Parameters map[string]string `json:"parameters"`
	}
	var sbr ServiceBindingRequest
	if err := json.NewDecoder(r.Body).Decode(&sbr); err != nil {
		return nil
	}
	return sbr.Parameters
}
