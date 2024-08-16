package common

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	. "terraform-provider-idmc/internal/utils"
)

// ClientResponse
// Basic details of a parsed api response.
type ClientResponse struct {
	*http.Response
	Body []byte
}

// NewClientResponse
// Wraps a provided http response.
func NewClientResponse(response *http.Response) (ClientResponse, error) {

	// Safely acquire the response body
	body := response.Body
	defer func() { _ = body.Close() }()

	// Read the body in its entirety so it can be used multiple times.
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return BadVal(err, ClientResponse{
			Response: response,
		})
	}

	return OkVal(ClientResponse{
		Response: response,
		Body:     bodyBytes,
	})

}

// IdmcClientResponse
// Currently aspirational, but intended to be a base struct for all IDMC api
// responses.
type IdmcClientResponse[DataType any, ErrorType any] struct {
	ClientResponse
	Success *DataType
	Failure *ErrorType
}

func NewIdmcClientResponse[DataType any, ErrorType any](rsp *http.Response) (IdmcClientResponse[DataType, ErrorType], error) {

	// Prepare the api response by wrapping the http response.
	clientResponse, err := NewClientResponse(rsp)
	response := IdmcClientResponse[DataType, ErrorType]{
		ClientResponse: clientResponse,
	}
	if err != nil {
		return response, err
	}

	// If we aren't dealing with a JSON payload, skip parsing.
	if !strings.Contains(rsp.Header.Get("Content-Type"), "json") {
		return response, nil
	}

	// Anything in the 2xx range is expected to be a success.
	if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
		var data DataType
		if err = json.Unmarshal(clientResponse.Body, &data); err != nil {
			return response, err
		}
		response.Success = &data
		return response, json.Unmarshal(clientResponse.Body, &response.Success)
	}

	// Handle any error messages, all expected to be the same format.
	return response, json.Unmarshal(response.Body, &response.Failure)

}
