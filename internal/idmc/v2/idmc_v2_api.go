package v2

import (
	"context"
	"net/http"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/utils"
)

type IdmcAdminV2Api struct {
	Client *ClientWithResponses
}

func NewIdmcAdminV2Api(baseUrl string, sessionId string, httpClient common.HttpRequestDoer) (*IdmcAdminV2Api, error) {

	apiClient, clientErr := NewClientWithResponses(baseUrl,
		common.WithHTTPClient(httpClient),
		common.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header["Accept"] = []string{"application/json"}
			req.Header["icSessionId"] = []string{sessionId}
			return nil
		}),
	)
	if clientErr != nil {
		return nil, clientErr
	}

	return utils.OkPtr(&IdmcAdminV2Api{
		Client: apiClient,
	})

}
