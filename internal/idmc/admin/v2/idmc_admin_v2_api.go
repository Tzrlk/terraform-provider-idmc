package v2

import (
	"context"
	"net/http"
)

type IdmcAdminV2Api struct {
	Client *ClientWithResponses
}

func NewIdmcAdminV2Api(baseUrl string, sessionId *string) (*IdmcAdminV2Api, error) {
	api := &IdmcAdminV2Api{}

	// Initialise the OpenAPI client.
	client, clientErr := NewClientWithResponses(baseUrl, func(httpClient *Client) error {

		// Inject needed headers.
		httpClient.RequestEditors = append(httpClient.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header["icSessionId"] = []string{*sessionId}

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
