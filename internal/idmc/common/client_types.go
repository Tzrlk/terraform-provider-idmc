package common

import (
	"net/http"
)

// ClientOption allows setting custom parameters during construction.
type ClientOption func(*ClientConfig) error

// HttpRequestDoer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client provides access to client configuration.
type Client interface {
	Config() *ClientConfig
}
