package idmc

import (
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"
)

type IdmcApi struct {
	V2 *v2.IdmcAdminV2Api
	V3 *v3.IdmcAdminV3Api
}

func NewIdmcApi(baseUrl string, sessionId string, httpClient common.HttpRequestDoer) (*IdmcApi, error) {

	idmcAdminV2Api, idmcAdminV2ApiErr := v2.NewIdmcAdminV2Api(baseUrl, sessionId, httpClient)
	if idmcAdminV2ApiErr != nil {
		return nil, idmcAdminV2ApiErr
	}

	idmcAdminV3Api, idmcAdminV3ApiErr := v3.NewIdmcAdminV3Api(baseUrl, sessionId, httpClient)
	if idmcAdminV3ApiErr != nil {
		return nil, idmcAdminV3ApiErr
	}

	return utils.OkPtr(&IdmcApi{
		V2: idmcAdminV2Api,
		V3: idmcAdminV3Api,
	})

}
