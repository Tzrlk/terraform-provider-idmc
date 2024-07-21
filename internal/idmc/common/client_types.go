package common

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function.
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// ResponseEditorFn  is the function signature for the RequestEditor callback function.
type ResponseEditorFn func(ctx context.Context, req *http.Response) error

// HttpRequestDoer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiResponse interface {
	Status() string
	StatusCode() int
	HttpResponse() *http.Response
	BodyData() []byte
}

type ClientConfig struct {

	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// Client provides access to client configuration.
type Client interface {
	Config() *ClientConfig
}

// ClientOption allows setting custom parameters during construction.
type ClientOption func(*ClientConfig) error

// NewClientConfig sets up a new ClientConfig with reasonable defaults.
func NewClientConfig(server string, opts ...ClientOption) (*ClientConfig, error) {
	config := ClientConfig{
		Server:         server,
		RequestEditors: make([]RequestEditorFn, 0),
	}

	// mutate client and add all optional params
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(config.Server, "/") {
		config.Server += "/"
	}

	// create httpClient, if not already present
	if config.Client == nil {
		config.Client = &http.Client{}
	}

	return &config, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(config *ClientConfig) error {
		config.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(config *ClientConfig) error {
		config.RequestEditors = append(config.RequestEditors, fn)
		return nil
	}
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(config *ClientConfig) error {
		newBaseURL, err := url.Parse(baseURL)
		if err == nil {
			config.Server = newBaseURL.String()
		}
		return err
	}
}

func (c *ClientConfig) ApplyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}
