package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/springernature/cf-route-service-sso-oauth/providers"
)

var (
	authUri  string
	tokenUri string
	clientID string
	clientS  string
)

func main() {
	log.SetOutput(os.Stdout)

	// Get initial ProviderData
	data := providers.InitProviderData()
	// Create Provider instance
	provider, err := providers.New(data)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// Sign in handler
	http.HandleFunc("/signin", provider.SignIn)
	// Callback handler
	http.HandleFunc("/callback", provider.Callback)

	// Default handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		// Check if already logged in
		cookie, _ := r.Cookie("auth")
		if cookie != nil {
			fmt.Fprint(w, "You are authenticated!")
		} else {
			fmt.Fprint(w, "<html>You are NOT authenticated. <a href=\"/signin\">Sign in</a> Host: "+r.Host+"</html>")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
