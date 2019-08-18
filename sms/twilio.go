package sms

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// ExceptionCode specific to twilio
type ExceptionCode int

const (
	baseURL  = "https://api.twilio.com/2010-04-01"
	videoURL = "https://video.twilio.com"

	clientTimeout = time.Second * 30
)

// The default http.Client that is used if none is specified
var defaultClient = &http.Client{
	Timeout: time.Second * 30,
}

type Twilio struct {
	AccountSid string
	AuthToken  string
	BaseURL    string
	VideoURL   string

	HTTPClient *http.Client

	APIKeySid    string
	APIKeySecret string
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int           `json:"status"`    // HTTP specific error code
	Message  string        `json:"message"`   // HTTP error message
	Code     ExceptionCode `json:"code"`      // Twilio specific error code
	MoreInfo string        `json:"more_info"` // Additional info from Twilio
}

// Create a new Twilio struct.
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return NewTwilioClientCustomHTTP(accountSid, authToken, nil)
}

// Create a new Twilio client, optionally using a custom http.Client
func NewTwilioClientCustomHTTP(accountSid, authToken string, HTTPClient *http.Client) *Twilio {
	if HTTPClient == nil {
		HTTPClient = defaultClient
	}

	return &Twilio{
		AccountSid: accountSid,
		AuthToken:  authToken,
		BaseURL:    baseURL,
		VideoURL:   videoURL,
		HTTPClient: HTTPClient,
	}
}

func (twilio *Twilio) WithAPIKey(apiKeySid string, apiKeySecret string) *Twilio {
	twilio.APIKeySid = apiKeySid
	twilio.APIKeySecret = apiKeySecret
	return twilio
}

func (twilio *Twilio) getBasicAuthCredentials() (string, string) {
	if twilio.APIKeySid != "" {
		return twilio.APIKeySid, twilio.APIKeySecret
	}

	return twilio.AccountSid, twilio.AuthToken
}

func (twilio *Twilio) post(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return twilio.do(req)
}

func (twilio *Twilio) get(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())

	return twilio.do(req)
}

func (twilio *Twilio) delete(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.getBasicAuthCredentials())

	return twilio.do(req)
}

func (twilio *Twilio) do(req *http.Request) (*http.Response, error) {
	client := twilio.HTTPClient
	if client == nil {
		client = defaultClient
	}

	return client.Do(req)
}

// Build path to a resource within the Twilio account
func (twilio *Twilio) buildUrl(resourcePath string) string {
	return twilio.BaseURL + "/" + path.Join("Accounts", twilio.AccountSid, resourcePath)
}
