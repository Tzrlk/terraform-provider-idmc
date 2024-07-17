package idmc

import (
	"terraform-provider-idmc/internal/idmc/admin"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/utils"
)

type IdmcApi struct {
	Admin *admin.IdmcAdminApi
}

func NewIdmcApi(baseApiUrl string, sessionId string, httpClient common.HttpRequestDoer) (*IdmcApi, error) {

	idmcAdminApi, idmcAdminApiErr := admin.NewIdmcAdminApi(baseApiUrl, sessionId, httpClient)
	if idmcAdminApiErr != nil {
		return nil, idmcAdminApiErr
	}

	return utils.OkPtr(&IdmcApi{
		Admin: idmcAdminApi,
	})

}
