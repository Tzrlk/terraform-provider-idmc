package common

import "net/http"

// ClientResponse
// Currently aspirational, but intended to be a base struct for all IDMC api
// responses.
type ClientResponse[S any, E any] struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *S
	JSON400      *E
	JSON401      *E
	JSON403      *E
	JSON404      *E
	JSON500      *E
	JSON502      *E
	JSON503      *E
}
