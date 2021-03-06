package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/springernature/oauth-route-service-broker/broker"
	"github.com/springernature/oauth-route-service-broker/handler"
	"github.com/springernature/oauth-route-service-broker/providers"
)

func main() {
	log.SetOutput(os.Stdout)
	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	// Set handlers for all enabled oauth providers
	for _, p := range providers.EnabledProviders() {
		// Sign in handler
		r.Get("/" + p.Name() + providers.SigninPath, p.SignIn)

		// Callback handler
		r.Handle("/" + p.Name() + providers.CallbackPath, handler.NewCallbackHandler(p))

		// Default handler
		// (i.e. Check if user is already authenticated. If yes, proxy to the app)
		r.HandleFunc("/" + p.Name(), handler.DefaultPathHandler)
	}

	/*
		Via the app:
		r.Host: sso-test-gerard.snpaas.eu
		Header: https://gerard-test.snpaas.eu/test/?bla=foo

		Direct on the service:
		r.URL.Path: /test
		r.URL.RequestURI: /test?bla=foo&animal=dog
		r.Host: sso-test-gerard.snpaas.eu
	*/

	// =============================
	// ====== Broker handlers ======
	// =============================

	// Service Catalog
	r.Get("/v2/catalog", broker.Catalog)

	// Service provisioning (Create, Delete and Status of a 'service instance' )
	r.Route("/v2/service_instances/{service_id}", broker.Provision)

	// Service binding
	r.Route("/v2/service_instances/{service_id}/service_bindings/{service_binding_id}", broker.Binding)

	// Default port number
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	// Start web server
	log.Fatal(http.ListenAndServe(":"+port, r))
}
