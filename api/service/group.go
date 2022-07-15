package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adrianosela/okta-gbac/api/service/payloads"
	"github.com/adrianosela/okta-gbac/utils/set"
	"github.com/gorilla/mux"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func (s *service) withGroupEndpoints() *service {
	s.router.Methods(http.MethodPost).Path("/group").Handler(s.auth(s.createGroupHandler))
	s.router.Methods(http.MethodPatch).Path("/group/{name}/add").Handler(s.auth(s.addToGroupHandler))
	s.router.Methods(http.MethodPatch).Path("/group/{name}/remove").Handler(s.auth(s.removeFromGroupHandler))
	s.router.Methods(http.MethodDelete).Path("/group/{name}").Handler(s.auth(s.deleteGroupHandler))
	return s
}

func (s *service) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	var pl *payloads.CreateGroupRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body not a valid CreateGroupRequest"))
		return
	}

	group := okta.Group{
		Profile: &okta.GroupProfile{
			Name:        pl.Name,
			Description: pl.Description,
			GroupProfileMap: map[string]interface{}{
				"api_managed": true,
				"owners":      set.NewSet(pl.Owners...).Add(authenticatedUser).Slice(),
			},
		},
	}

	// https://pkg.go.dev/github.com/okta/okta-sdk-golang/v2/okta#GroupResource.CreateGroup
	createdGroup, _, err := s.okta.Group.CreateGroup(context.Background(), group)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	gbytes, err := json.Marshal(&createdGroup)
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

// add members and owners
func (s *service) addToGroupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not implemented"))
	return
}

// remove members and owners
func (s *service) removeFromGroupHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not implemented"))
	return
}

func (s *service) deleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	gname := mux.Vars(r)["name"]
	if gname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No group name in request URL"))
		return
	}

	// https://pkg.go.dev/github.com/okta/okta-sdk-golang/v2/okta#GroupResource.ListGroups
	qp := query.NewQueryParams(query.WithSearch(
		fmt.Sprintf("profile.api_managed eq true and profile.name eq \"%s\"", gname),
	))
	groups, resp, err := s.okta.Group.ListGroups(context.Background(), qp)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	if len(groups) == 0 { // group does not exist
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// okta seems to disallow duplicate names so this should never happen.
	// still worth catching just in case.
	if len(groups) > 1 {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(fmt.Sprintf("Group \"%s\" is not unique. Contact your administrator", gname)))
		return
	}

	group := groups[0]
	if !isOwner(group, authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("You can not delete a group you are not an owner of!"))
		return
	}

	resp, err = s.okta.Group.DeleteGroup(context.Background(), group.Id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An unexpected error occured"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
