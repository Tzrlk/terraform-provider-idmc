package common

import (
	"context"
	"net/http"
	"strings"
)

type ClientConfig struct {

	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A collection of callbacks for modifying requests and responses handled
	// by this client configuration.
	Editors ClientConfigEditor
}

// NewClientConfig sets up a new ClientConfig with reasonable defaults
func NewClientConfig(server string, opts ...ClientOption) (*ClientConfig, error) {
	config := ClientConfig{
		Server: server,
		Editors: ClientConfigEditor{
			RequestEditors:     make([]RequestEditorFn, 0),
			ResponseEditors:    make([]ResponseEditorFn, 0),
			ApiResponseEditors: make([]ApiResponseEditorFn, 0),
		},
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

func (c *ClientConfig) HandleRequest(
	ctx context.Context,
	editors []ClientConfigEditor,
	create func() (*http.Request, error),
) (*http.Response, error) {

	// Generate the API request
	req, err := create()
	if err != nil {
		return nil, err
	}

	// Enrich request with
	req = req.WithContext(ctx)

	// Merge editors for this request in prep for usage.
	editor := c.Editors.Merge(editors...)

	// Apply request editors
	if err := editor.EditHttpRequest(ctx, req); err != nil {
		return nil, err
	}

	// Perform the request
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Apply response editors
	if err := editor.EditHttpResponse(ctx, res); err != nil {
		return nil, err
	}

	// Return the response
	return res, nil
}
