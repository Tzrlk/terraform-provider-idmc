package v3

import (
	"context"
	"net/http"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/utils"
)

type IdmcAdminV3Api struct {
	Client *ClientWithResponses
}

func NewIdmcAdminV3Api(baseUrl string, sessionId string, httpClient common.HttpRequestDoer) (*IdmcAdminV3Api, error) {

	// Initialise the OpenAPI client.
	apiClient, clientErr := NewClientWithResponses(baseUrl,
		common.WithHTTPClient(httpClient),
		common.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header["Accept"] = []string{"application/json"}
			req.Header["INFA-SESSION-ID"] = []string{sessionId}
			return nil
		}),
	)
	if clientErr != nil {
		return nil, clientErr
	}

	return utils.OkPtr(&IdmcAdminV3Api{
		Client: apiClient,
	})

}
