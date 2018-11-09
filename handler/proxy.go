package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/springernature/oauth-route-service-broker/htmltemplate"
	"github.com/springernature/oauth-route-service-broker/providers"
	"github.com/springernature/oauth-route-service-broker/token"
)

func DefaultPathHandler(w http.ResponseWriter, r *http.Request) {
	// CF sends the user requested Url in this header
	forwardUrl := r.Header.Get("X-Cf-Forwarded-Url")
	if forwardUrl == "" {
		fmt.Fprint(w, htmltemplate.NoForwardUrlErr)
		return
	}

	// ================ PROXY ================ //
	// If authentication token (JWT) in cookie is still valid, proxy the request to the target app
	jwtcookie, _ := r.Cookie("CfSsoJwt")
	if jwtcookie != nil && token.ValidJwt(jwtcookie.Value) {
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
		return
	}

	// ==== NO SUCCESSFUL AUTHENTICATION ==== //
	// If authentication is not valid, do a redirect to the signin path on the sso service
	// Compare if sso service and the app run on the same apex domain
	// First get the apex from the app
	url, _ := url.Parse(forwardUrl)
	appApex := apexDomain(url.Host)
	// Select the sso service uri that matches on apex level with the app apex
	ssoUri := ssoUriSelector(appApex)
	if ssoUri == "" {
		fmt.Fprintf(w, htmltemplate.NotSameApexErr, appApex)
		return
	}
	// We want to know where to redirect to after signin in.
	// This is stored in a temporary cookie
	c := &http.Cookie{
		Name:   "CfSsoReqOrigin",
		Value:  forwardUrl,
		Domain: appApex,
		Path:   "/",
		MaxAge: 300,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "https://"+ssoUri+r.URL.Path+providers.SigninPath, 302)
}

func director(req *http.Request) {
	forwardedURL := req.Header.Get("X-Cf-Forwarded-Url")

	url, err := url.Parse(forwardedURL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	req.URL = url
	req.Host = url.Host
}

func apexDomain(d string) string {
	s := strings.Split(d, ".")
	if len(s) < 2 {
		return ""
	}
	return s[len(s)-2] + "." + s[len(s)-1]
}

// Checks if the app apex matches with one of the apexs from the domains where the
// sso service runs on. If there is a match, return which sso service domain to use.
func ssoUriSelector(appApex string) string {
	// Get the domain(s) where this sso service runs on
	vcap := os.Getenv("VCAP_APPLICATION")
	type VcapApplication struct {
		Uris []string `json:"uris"`
	}
	var va VcapApplication
	err := json.Unmarshal([]byte(vcap), &va)
	if err != nil {
		log.Printf("Error while parsing VCAP_APPLICATION in 'sameApex() func: %v\n", err)
		return ""
	}
	// Check if there is a match between app (forwardUrl) and sso service domains
	// Loop over all uris for sso service
	for _, d := range va.Uris {
		apex := apexDomain(d)
		if apex == appApex {
			// Match! Return the sso service uri where the match was found
			return d
		}
	}
	return ""
}
