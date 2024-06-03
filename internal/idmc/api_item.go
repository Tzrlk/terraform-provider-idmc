package idmc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ApiItem[I any, O any] struct {
	Api    *Api
	Method string
	Path   string
}

func (fn *ApiItem[I, O]) BuildHeaders() http.Header {

	headers := http.Header{
		"Accept": {"application/json"},
	}

	// Add session id and base url to headers if in a session.
	if fn.Api.SessionId != "" {
		if strings.Contains(fn.Path, "/v2/") {
			headers.Set("icSessionId", fn.Api.SessionId)
			headers.Set("serverUrl", fn.Api.BaseUrl)

		} else if strings.Contains(fn.Path, "/v3/") {
			headers.Set("INFA-SESSION-ID", fn.Api.SessionId)
			headers.Set("baseApiUrl", fn.Api.BaseUrl)
		}
	}

	return headers
}

func (fn *ApiItem[I, O]) SerialiseRequestBody(
	diag *diag.Diagnostics,
	requestBody *I,
) *bytes.Buffer {

	if requestBody == nil {
		return nil
	}

	requestBodyReader := new(bytes.Buffer)
	err := json.NewEncoder(requestBodyReader).Encode(requestBody)
	if err != nil {
		diag.AddError(
			"JSON request serialisation error",
			fmt.Sprintf("%s", err),
		)
		return nil
	}

	return requestBodyReader
}

func (fn *ApiItem[I, O]) BuildRequestUrl(
	pathParams []any,
	queryParams map[string]string,
) string {
	requestUrl := fn.Api.BaseUrl + "/" + fn.Path

	if pathParams != nil && strings.Contains(requestUrl, "%s") {
		requestUrl = fmt.Sprintf(requestUrl, pathParams...)
	}

	if queryParams != nil {
		requestUrl += "?"
		for key, val := range queryParams {
			requestUrl += key + "=" + val
		}
	}

	return requestUrl
}

func (fn *ApiItem[I, O]) DeserialiseResponseBody(
	diag *diag.Diagnostics,
	resp *http.Response,
) *O {

	// Ensure the response body is closed after use.
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			diag.AddError(
				"Request Body Closing Error",
				fmt.Sprintf("Encountered an issue closing a response body: %s", closeErr),
			)
		}
	}(resp.Body)

	// Deserialise the response body.
	var responseBody *O
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		diag.AddError(
			"JSON response parsing error",
			fmt.Sprintf("%s", err),
		)
		return nil
	}

	return responseBody

}

func (fn *ApiItem[I, O]) Call(
	diag *diag.Diagnostics,
	pathParams []any,
	queryParams map[string]string,
	requestBody *I,
) *O {

	// Ensure that the request body is serialised to JSON if provided.
	requestBodyReader := fn.SerialiseRequestBody(diag, requestBody)
	if diag.HasError() {
		return nil
	}

	// Set up standard api request headers
	headers := fn.BuildHeaders()
	if requestBodyReader != nil {
		headers.Set("Content-Type", "application/json")
	}

	// Set up the request.
	requestUrl      := fn.BuildRequestUrl(pathParams, queryParams)
	request, reqErr := http.NewRequest(fn.Method, requestUrl, requestBodyReader)
	if reqErr != nil {
		diag.AddError(
			"Request Creation Error",
			fmt.Sprintf("Unable to create http client request: %s", reqErr),
		)
		return nil
	}

	// Add our previously set-up headers to the request.
	request.Header = headers

	// Perform the request
	response, httpErr := fn.Api.Client.Do(request)
	if httpErr != nil {
		diag.AddError(
			"Client Error",
			fmt.Sprintf("Unable to complete api request: %s", httpErr),
		)
		return nil
	}

	// Handle the response.
	return fn.DeserialiseResponseBody(diag, response)

}
