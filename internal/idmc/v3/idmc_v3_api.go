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

func NewIdmcAdminV3Api(baseUrl string, sessionId *string, opts ...common.ClientOption) (*IdmcAdminV3Api, error) {

	// Add a request editor to apply the needed api headers on all requests.
	opts = append(opts, common.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header["Accept"] = []string{"application/json"}
		if sessionId != nil {
			req.Header["INFA-SESSION-ID"] = []string{*sessionId}
		}
		return nil
	}))

	apiClient, clientErr := NewClientWithResponses(baseUrl, opts...)
	if clientErr != nil {
		return nil, clientErr
	}

	return utils.OkPtr(&IdmcAdminV3Api{
		Client: apiClient,
	})

}
