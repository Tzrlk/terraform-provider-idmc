package v3

import (
	"context"
	"net/http"
)

type IdmcAdminV3Api struct {
	Client *ClientWithResponses
}

func NewIdmcAdminV3Api(baseUrl string, sessionId *string) (*IdmcAdminV3Api, error) {
	api := &IdmcAdminV3Api{}

	// Initialise the OpenAPI client.
	client, clientErr := NewClientWithResponses(baseUrl, func(httpClient *Client) error {

		// Inject needed headers.
		httpClient.RequestEditors = append(httpClient.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header["INFA-SESSION-ID"] = []string{*sessionId}

			return nil
		})

		return nil
	})
	if clientErr != nil {
		return nil, clientErr
	}

	api.Client = client
	return api, nil

}
