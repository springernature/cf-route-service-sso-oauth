package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/springernature/cf-route-service-sso-oauth/handler"
	"github.com/springernature/cf-route-service-sso-oauth/providers"
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
	h := &handler.CallbackHandler{
		Redeem:   provider.Redeem,
		GetEmail: provider.GetEmail,
		Filter:   provider.Filter,
	}
	http.Handle("/callback", handler.NewCallbackHandler(h))

	// Default handler
	// (i.e. Check if user is already authenticated. If yes, proxy to the app)
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

	// Default port number
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	// Start web server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
