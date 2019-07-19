package util

import (
	"net/http"
	"time"
)

// DefaultHTTPClient with a timeout. Default implementation does not a timeout
var DefaultHTTPClient = &http.Client{
	Timeout: time.Second * 10,
}
