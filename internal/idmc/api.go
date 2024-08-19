package idmc

import (
	"context"
	"fmt"
	"net/http"

	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"
)

type IdmcApi struct {
	V2 v2.IdmcAdminV2Api
	V3 v3.IdmcAdminV3Api
}

const msgInitFailed = "unable to initialise api for %s: %v"

func NewIdmcApi(baseUrl string, opts ...common.ClientOption) (*IdmcApi, error) {

	// Add a request editor to add common api headers on all requests.
	opts = append(opts, common.WithRequestEditorFn(func(ctx context.Context, cfg *common.ClientConfig, req *http.Request) error {
		req.Header["Accept"] = []string{"application/json"}
		return nil
	}))

	// Set up the common client configuration.
	config, err := common.NewClientConfig(baseUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf(msgInitFailed, baseUrl,
			fmt.Errorf("failed to create common client config: %v", err))
	}

	// Initialise the v2 api.
	apiV2 := v2.NewIdmcAdminV2Api(config)

	// Initialise the v3 api.
	apiV3 := v3.NewIdmcAdminV3Api(config)

	// Construct the api now everything is ok.
	return utils.OkPtr(&IdmcApi{
		V2: apiV2,
		V3: apiV3,
	})

}
