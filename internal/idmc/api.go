package idmc

import (
	"fmt"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"
)

type IdmcApi struct {
	V2 v2.IdmcAdminV2Api
	V3 v3.IdmcAdminV3Api
}

func NewIdmcApi(baseUrl string, sessionId string, opts ...common.ClientOption) (*IdmcApi, error) {

	apiV2, err := v2.NewIdmcAdminV2Api(baseUrl, &sessionId, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise api for %s: %v", baseUrl, err)
	}

	apiV3, err := v3.NewIdmcAdminV3Api(baseUrl, &sessionId, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise api for %s: %v", baseUrl, err)
	}

	return utils.OkPtr(&IdmcApi{
		V2: *apiV2,
		V3: *apiV3,
	})

}
