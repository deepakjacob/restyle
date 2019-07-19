package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/deepakjacob/restyle/util"
)

type httpClient struct {
	Client *http.Client
}

// Client client for contactacting google service
var Client = &httpClient{
	Client: util.DefaultHTTPClient,
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
