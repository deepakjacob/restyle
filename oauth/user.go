package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deepakjacob/restyle/util"
)

// GoogleUser struct google combined scopes profile and email returns
type GoogleUser struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_name"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	NickName      string `json:"nick_name"`
	OrgID         string `json:"org_id"`
	UserID        string `json:"user_id"`
}

// HTTPClient for http operations
type HTTPClient interface {
	Get(string) (*GoogleUser, error)
}

func (h *httpClient) Get(url string) (*GoogleUser, error) {
	response, err := h.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	var user GoogleUser
	if err := json.Unmarshal([]byte(contents), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user json to struct %v", err.Error())
	}
	return &user, nil
}

type httpClient struct {
	Client *http.Client
}

// Client client for contactacting google service
var Client = &httpClient{
	Client: util.DefaultHTTPClient,
}
