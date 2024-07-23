package common

import "net/url"

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
