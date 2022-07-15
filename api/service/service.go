package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// Config represents configuration for the service
type Config struct {
	OktaNamespace           string
	OktaAPIToken            string
	OktaAPIRequestTimeout   int64
	OktaAPIRateLimitRetries int32
}

type service struct {
	router *mux.Router
	okta   *okta.Client
}

// New returns the handler for a new service
func New(c Config) (http.Handler, error) {
	// TODO: validate config object

	_, oktaClient, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(fmt.Sprintf("https://%s.com", c.OktaNamespace)),
		okta.WithToken(c.OktaAPIToken),
		okta.WithRequestTimeout(c.OktaAPIRequestTimeout),
		okta.WithRateLimitMaxRetries(c.OktaAPIRateLimitRetries),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Okta Developer API client: %s", err)
	}

	service := &service{
		router: mux.NewRouter(),
		okta:   oktaClient,
	}

	return service.
		withDebugEndpoints().
		withGroupEndpoints().
		withUserEndpoints().
		router, nil
}
