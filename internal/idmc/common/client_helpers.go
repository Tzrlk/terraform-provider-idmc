package common

import (
	"context"
	"net/http"
)

// ClientOption allows setting custom parameters during construction.
type ClientOption func(*ClientConfig) error

// ClientOptions a collection of ClientOption items.
type ClientOptions []ClientOption

// WithRequestEditorFn see common.WithRequestEditorFn.
func (c ClientOptions) WithRequestEditorFn(fn RequestEditorFn) ClientOptions {
	return append(c, WithRequestEditorFn(fn))
}

// WithRequestHeader see common.WithRequestHeader.
func (c ClientOptions) WithRequestHeader(name string, values ...string) ClientOptions {
	return c.WithRequestEditorFn(WithRequestHeader(name, values...))
}

////////////////////////////////////////////////////////////////////////////////

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
		config.Editors.RequestEditors = append(config.Editors.RequestEditors, fn)
		return nil
	}
}

// WithResponseEditorFn allows setting up a callback function, which will be
// called right after receiving the response. This can be used to mutate the response.
func WithResponseEditorFn(fn ResponseEditorFn) ClientOption {
	return func(config *ClientConfig) error {
		config.Editors.ResponseEditors = append(config.Editors.ResponseEditors, fn)
		return nil
	}
}

// WithApiResponseEditorFn allows setting up a callback function, which will be
// called right after parsing the response. This can be used to mutate the response.
func WithApiResponseEditorFn(fn ApiResponseEditorFn) ClientOption {
	return func(config *ClientConfig) error {
		config.Editors.ApiResponseEditors = append(config.Editors.ApiResponseEditors, fn)
		return nil
	}
}

// WithRequestHeader adds the specified header to any requests with the provided
// values.
func WithRequestHeader(name string, values ...string) RequestEditorFn {
	return func(ctx context.Context, cfg *ClientConfig, req *http.Request) error {
		req.Header[name] = values
		return nil
	}
}

// WithSessionHeader injects the current session id into requests
// as a header with the provided name.
func WithSessionHeader(name string) RequestEditorFn {
	return func(ctx context.Context, cfg *ClientConfig, req *http.Request) error {
		req.Header[name] = []string{cfg.SessionId}
		return nil
	}
}
