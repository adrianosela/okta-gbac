package service

import (
	"context"
	"net/http"
)

var (
	authenticatedUserContextKey = "authenticated-user"
)

// auth wraps a handler function with authenicated
func (s *service) auth(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// FIXME: Get JWT from "Authorization" header, validate it, inject user into context

		username := r.Header.Get("MOCK_AUTHENTICATED_USER")
		if username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("No user in \"MOCK_AUTHENTICATED_USER\" header"))
			return
		}

		// run handler with auth values in context
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), authenticatedUserContextKey, username)))
	})
}

// getAuthenticatedUser returns the authenticated user in the context object
func getAuthenticatedUser(r *http.Request) string {
	return r.Context().Value(authenticatedUserContextKey).(string)
}
