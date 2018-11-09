package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	if !sameApex(apexDomain(r.Host), forwardUrl) {
		fmt.Fprintf(w, htmltemplate.NotSameApexErr, apexDomain(r.Host))
		return
	}
	// We want to know where to redirect to after signin in.
	// This is stored in a temporary cookie
	c := &http.Cookie{
		Name:   "CfSsoReqOrigin",
		Value:  forwardUrl,
		Domain: apexDomain(r.Host),
		Path:   "/",
		MaxAge: 300,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "https://"+r.Host+r.URL.Path+providers.SigninPath, 302)
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

func sameApex(d1 string, f string) bool {
	// Parse domain from forwardUrl
	url, err := url.Parse(f)
	if err != nil {
		log.Printf("Error while parsing forwardUrl in 'sameApex() func: %v\n", err)
		return false
	}
	// Get apex domain from forwardUrl hostname
	d2 := apexDomain(url.Host)
	if d1 == d2 {
		return true
	}
	return false
}
