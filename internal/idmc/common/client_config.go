package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type ClientConfig struct {

	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// The session id acquired from a successful login operation. Should be
	// passed along on every request if not blank.
	SessionId string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A collection of callbacks for modifying requests and responses handled
	// by this client configuration.
	Editors ClientConfigEditor
}

// NewClientConfig sets up a new ClientConfig with reasonable defaults.
func NewClientConfig(server string, opts ...ClientOption) (*ClientConfig, error) {
	config := ClientConfig{
		Server:  server,
		Editors: ClientConfigEditor{},
	}

	// mutate client and add all optional params.
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	// ensure the server URL always has a trailing slash.
	if !strings.HasSuffix(config.Server, "/") {
		config.Server += "/"
	}

	// create httpClient, if not already present.
	if config.Client == nil {
		config.Client = &http.Client{}
	}

	return &config, nil
}

func (c *ClientConfig) SetServer(server string) error {

	newBaseURL, err := url.Parse(server)
	if err != nil {
		return fmt.Errorf("unable to set server url to %s: %v", server, err)
	}

	newBaseURL.Scheme = "https" // Ensure we're using https at all times.
	c.Server = newBaseURL.String()

	// ensure the server URL always has a trailing slash.
	if !strings.HasSuffix(c.Server, "/") {
		c.Server += "/"
	}

	return nil
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

	// Enrich request with the current context.
	req = req.WithContext(ctx)

	// Merge editors for this request in prep for usage.
	editor := c.Editors.Merge(editors...)

	// Apply request editors
	if err := editor.EditHttpRequest(ctx, c, req); err != nil {
		return nil, err
	}

	// Perform the request.
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Apply response editors.
	if err := editor.EditHttpResponse(ctx, c, res); err != nil {
		return nil, err
	}

	// Return the response.
	return res, nil
}

func HandleAction[T any](
	client *ClientConfig,
	ctx context.Context,
	editors []ClientConfigEditor,
	parser func(rsp *http.Response) (*T, error),
	action func(ctx context.Context) (*http.Response, error),
) (*T, error) {

	// Perform the provided action.
	rsp, err := action(ctx)
	if err != nil {
		return nil, err
	}

	// Parse the response payload.
	apiRes, err := parser(rsp)
	if err != nil {
		return nil, err
	}

	// Sketchy reflection because we can't have nice things.
	apiResVal := reflect.ValueOf(apiRes).Elem()
	clientRspVal := apiResVal.FieldByName("ClientResponse")

	clientResp, ok := clientRspVal.Interface().(ClientResponse)
	if !ok {
		return nil, errors.New("failed to convert apiRes field to ClientResponse")
	}

	// Apply API response editors.
	editor := client.Editors.Merge(editors...)
	if err := editor.EditApiResponse(ctx, client, &clientResp); err != nil {
		return nil, err
	}

	clientRspVal.Set(reflect.ValueOf(clientResp))

	// Success.
	return apiRes, nil

}
