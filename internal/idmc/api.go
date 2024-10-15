package idmc

import (
	"fmt"

	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"

	. "github.com/samber/mo"
	. "terraform-provider-idmc/internal/idmc/common"
	. "terraform-provider-idmc/internal/utils"
)

type IdmcApi struct {
	V2 v2.IdmcAdminV2Api
	V3 v3.IdmcAdminV3Api
}

const msgInitFailed = "unable to initialise api for %s: %v"

func NewIdmcApi(baseUrl string, opts ClientOptions) Result[*IdmcApi] {

	// Set up the common client configuration.
	config := NewClientConfig(baseUrl, opts...)

	// Ensure any errors thrown are appropriately handled.
	config = config.MapErr(func(err error) (*ClientConfig, error) {
		return nil, fmt.Errorf(msgInitFailed, baseUrl,
			fmt.Errorf("failed to create common client config: %v", err))
	})

	// Construct the api now everything is ok.
	return MapResultOk(config, func(config *ClientConfig) Result[*IdmcApi] {
		return Ok(&IdmcApi{
			V2: v2.NewIdmcAdminV2Api(config),
			V3: v3.NewIdmcAdminV3Api(config),
		})
	})

}
