package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func (s *service) withUserEndpoints() *service {
	s.router.Methods(http.MethodGet).Path("/user/{username}/member").HandlerFunc(s.userMemberGroupsHandler)
	s.router.Methods(http.MethodGet).Path("/user/{username}/owner").HandlerFunc(s.userOwnerGroupsHandler)
	return s
}

func (s *service) userMemberGroupsHandler(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]
	if uname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No username in request URL"))
		return
	}

	// https://pkg.go.dev/github.com/okta/okta-sdk-golang/v2/okta#UserResource.ListUserGroups
	groups, _, err := s.okta.User.ListUserGroups(context.Background(), uname)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	names := []string{}
	for _, g := range groups {
		names = append(names, g.Profile.Name)
	}

	gbytes, err := json.Marshal(&names)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(gbytes)
	return
}

func (s *service) userOwnerGroupsHandler(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]
	if uname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No username in request URL"))
		return
	}

	// https://pkg.go.dev/github.com/okta/okta-sdk-golang/v2/okta#GroupResource.ListGroups
	qp := query.NewQueryParams(query.WithSearch("profile.api_managed eq true"))
	groups, _, err := s.okta.Group.ListGroups(context.Background(), qp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	// unfortunately we can't directly query for the "owners" array
	// containing an element, so we must filter out the groups ourselves
	names := []string{}
	for _, g := range groups {
		if isOwner(g, uname) {
			names = append(names, g.Profile.Name)
		}
	}

	gbytes, err := json.Marshal(&names)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(gbytes)
	return
}
