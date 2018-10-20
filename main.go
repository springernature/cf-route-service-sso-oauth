package main

import (
	"log"
	"net/http"
	"os"

	"github.com/springernature/oauth-route-service-broker/handler"
	"github.com/springernature/oauth-route-service-broker/providers"
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
	http.HandleFunc(providers.SigninPath, provider.SignIn)

	// Callback handler
	http.Handle(providers.CallbackPath, handler.NewCallbackHandler(provider))

	// Default handler
	// (i.e. Check if user is already authenticated. If yes, proxy to the app)
	http.HandleFunc("/", handler.DefaultPathHandler)

	/*
		Via the app:
		r.Host: sso-test-gerard.snpaas.eu
		Header: https://gerard-test.snpaas.eu/test/?bla=foo

		Direct on the service:
		r.URL.Path: /test
		r.URL.RequestURI: /test?bla=foo&animal=dog
		r.Host: sso-test-gerard.snpaas.eu
	*/

	// Default port number
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	// Start web server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
