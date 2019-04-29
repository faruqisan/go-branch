package httpclient

import (
	"net/http"
	"time"
)

var (
	defaultClient *Client
)

func init() {
	defaultClient = NewClient(WithHTTPTimeout(defaultHTTPTimeout * time.Second))
}

// Do executes the given http request and returns the http response.
// It uses default shared http client with default timeout value.
func Do(req *http.Request) (*http.Response, error) {
	return defaultClient.Do(req)
}

// DoJSON executes the given http request and unmarshall the response body
// into the given `data`
// It uses default shared http client with default timeout value.
func DoJSON(req *http.Request, data interface{}) (*http.Response, error) {
	return defaultClient.DoJSON(req, data)
}
