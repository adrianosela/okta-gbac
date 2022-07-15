package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func unmarshalRequestBody(r *http.Request, intf interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bodyBytes, intf); err != nil {
		return err
	}
	return nil
}

func isOwner(g *okta.Group, uname string) bool {
	if g != nil {
		for _, owner := range g.Profile.GroupProfileMap["owners"].([]interface{}) {
			if owner.(string) == uname {
				return true
			}
		}
	}
	return false
}
