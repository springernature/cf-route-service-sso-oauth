package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/springernature/oauth-route-service-broker/htmltemplate"
	"github.com/springernature/oauth-route-service-broker/providers"
	"github.com/springernature/oauth-route-service-broker/token"
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
	// Error returned means the user is not authenticated based on this filter
	attr, err := ch.Provider.Filter(b)
	if err != nil {
		fmt.Fprintf(w, htmltemplate.FilterErr, err)
		return
	}

	// ============= AUTHENTICATED USER ============= //
	// Issue a new JWT
	jwt, err := token.NewJwt(username, attr)
	if err != nil {
		fmt.Fprintf(w, htmltemplate.JwtIssueErr, err)
		return
	}

	// Fetch the URL of the app (origin)
	// We'll use this later to redirect back to the app
	orgCookie, err := r.Cookie("CfSsoReqOrigin")
	if err != nil {
		fmt.Fprintf(w, htmltemplate.NoOriginCookieErr, err)
		return
	}
	_, err = url.Parse(orgCookie.Value)
	if err != nil {
		fmt.Fprintf(w, htmltemplate.NoValidOriginErr, err)
		return
	}

	// Set the JWT cookie. Based on this cookie this service will now know
	// this user is already authenticated.
	jwtCookie := &http.Cookie{
		Name:   "CfSsoJwt",
		Value:  jwt,
		Domain: apexDomain(r.Host),
		Path:   "/",
	}
	http.SetCookie(w, jwtCookie)

	// Unset the origin cookie as we don't need it anymore
	// as we'll be redirecting back to the origin shortly after
	orgUrl := orgCookie.Value
	orgCookie.MaxAge = -1
	orgCookie.Domain = apexDomain(r.Host)
	orgCookie.Path = "/"
	http.SetCookie(w, orgCookie)
	// Redirect
	http.Redirect(w, r, orgUrl, 302)
}
