package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adrianosela/okta-gbac/api/service"
)

func main() {
	config := service.Config{
		OktaNamespace:           os.Getenv("OKTA_NAMESPACE"),
		OktaAPIToken:            os.Getenv("OKTA_API_TOKEN"),
		OktaAPIRequestTimeout:   30,
		OktaAPIRateLimitRetries: 3,
	}

	handler, err := service.New(config)
	if err != nil {
		log.Fatalf("Failed to initialize service: %s", err)
	}

	if err = http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to serve HTTP: %s", err)
	}
}
