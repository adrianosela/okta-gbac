package service

import (
	"fmt"
	"net/http"
)

func (s *service) withDebugEndpoints() *service {
	s.router.Methods(http.MethodGet).Path("/healthcheck").HandlerFunc(s.healthcheckHandler)
	s.router.Methods(http.MethodGet).Path("/authcheck").Handler(s.auth(s.authcheckHandler))
	return s
}

func (s *service) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm alive!"))
	return
}

func (s *service) authcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User \"%s\" is authenticated!", getAuthenticatedUser(r))))
	return
}
