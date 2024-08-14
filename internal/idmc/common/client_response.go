package common

import "net/http"

// ClientResponse
// Basic details of a parsed api response.
type ClientResponse struct {
	*http.Response
	Body []byte
}

// IdmcClientResponse
// Currently aspirational, but intended to be a base struct for all IDMC api
// responses.
type IdmcClientResponse[Dat any, Err any] struct {
	ClientResponse
	Success *Dat
	Failure *Err
}
