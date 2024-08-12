package v2

import (
	"context"
	"net/http"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/utils"
)

//go:generate -command oapi-codegen go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
//go:generate oapi-codegen -config ./codegen.yml ./openapi.yml

type IdmcAdminV2Api struct {
	Client ClientWithResponses
}

func NewIdmcAdminV2Api(baseUrl string, sessionId *string, opts ...common.ClientOption) (*IdmcAdminV2Api, error) {

	// Add a request editor to apply the needed api headers on all requests.
	opts = append(opts, common.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header["Accept"] = []string{"application/json"}
		if sessionId != nil {
			req.Header["icSessionId"] = []string{*sessionId}
		}
		return nil
	}))

	apiClient, clientErr := NewClientWithResponses(baseUrl, opts...)
	if clientErr != nil {
		return nil, clientErr
	}

	return utils.OkPtr(&IdmcAdminV2Api{
		Client: *apiClient,
	})

}
