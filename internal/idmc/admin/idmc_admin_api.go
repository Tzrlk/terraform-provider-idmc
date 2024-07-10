package admin

import (
	"terraform-provider-idmc/internal/idmc/admin/v2"
	"terraform-provider-idmc/internal/idmc/admin/v3"
)

type IdmcAdminApi struct {
	V2 *v2.IdmcAdminV2Api
	V3 *v3.IdmcAdminV3Api
}

func NewIdmcAdminApi(baseUrl string, sessionId *string) (*IdmcAdminApi, error) {
	api := &IdmcAdminApi{}

	idmcAdminV2Api, idmcAdminV2ApiErr := v2.NewIdmcAdminV2Api(baseUrl, sessionId)
	if idmcAdminV2ApiErr != nil {
		return nil, idmcAdminV2ApiErr
	}

	idmcAdminV3Api, idmcAdminV3ApiErr := v3.NewIdmcAdminV3Api(baseUrl, sessionId)
	if idmcAdminV3ApiErr != nil {
		return nil, idmcAdminV3ApiErr
	}

	api.V2 = idmcAdminV2Api
	api.V3 = idmcAdminV3Api
	return api, nil

}
