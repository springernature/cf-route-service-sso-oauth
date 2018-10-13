package handler

import (
	"fmt"
	"net/http"
)

type CallbackHandler struct {
	Redeem   func(*http.Request) ([]byte, error)
	GetEmail func([]byte) (string, error)
	Filter   func([]byte) (bool, error)
}

func NewCallbackHandler(ch *CallbackHandler) http.Handler {
	return ch
}

func (ch *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Redeem oauth tokens if valid authorization code is provided
	b, err := ch.Redeem(r)
	if err != nil {
		fmt.Fprintf(w, "Error while redeeming authorization code: %v", err)
		return
	}
	// Get the users email from the oauth payload
	username, err := ch.GetEmail(b)
	if err != nil {
		fmt.Fprintf(w, "Error while reading user email from oauth payload: %v", err)
		return
	}
	// Apply additional provider specifc filters (e.g. Is user member of group or team?)
	access, err := ch.Filter(b)
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
