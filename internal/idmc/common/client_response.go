package common

import "net/http"

// ClientResponse
// Basic details of a parsed api response.
type ClientResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// IdmcClientResponse
// Currently aspirational, but intended to be a base struct for all IDMC api
// responses.
type IdmcClientResponse[Dat any, Err any] struct {
	ClientResponse
	JSON200 *Dat
	JSON400 *Err
	JSON401 *Err
	JSON403 *Err
	JSON404 *Err
	JSON500 *Err
	JSON502 *Err
	JSON503 *Err
}
