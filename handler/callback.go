package handler

import (
	"fmt"
	"net/http"

	"github.com/springernature/cf-route-service-sso-oauth/providers"
)

type CallbackHandler struct {
	Provider providers.Provider
}

func NewCallbackHandler(p providers.Provider) http.Handler {
	return &CallbackHandler{
		Provider: p,
	}
}

func (ch *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Redeem oauth tokens if valid authorization code is provided
	b, err := ch.Provider.Redeem(r)
	if err != nil {
		fmt.Fprintf(w, "Error while redeeming authorization code: %v", err)
		return
	}
	// Get the users email from the oauth payload
	username, err := ch.Provider.GetEmail(b)
	if err != nil {
		fmt.Fprintf(w, "Error while reading user email from oauth payload: %v", err)
		return
	}
	// Apply additional provider specifc filters (e.g. Is user member of group or team?)
	access, err := ch.Provider.Filter(b)
	if err != nil {
		fmt.Fprintf(w, "Error while applying additional provider specific authorization filter: %v", err)
		return
	}
	if !access {
		fmt.Fprint(w, "Unauthorized!")
		return
	}

	// Redirect back to app

	// Return username
	fmt.Fprintf(w, "Username: %v", username)
}
