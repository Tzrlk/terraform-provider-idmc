package v2

import (
	"terraform-provider-idmc/internal/idmc/common"
)

//go:generate -command oapi-codegen go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
//go:generate oapi-codegen -config ./codegen.yml -generate types -o ./types.gen.go ./openapi.yml
//go:generate oapi-codegen -config ./codegen.yml -generate client -o ./client.gen.go ./openapi.yml

type IdmcAdminV2Api struct {
	Client ClientWithResponses
}

func NewIdmcAdminV2Api(config *common.ClientConfig) IdmcAdminV2Api {

	editors := common.ClientConfigEditor{
		RequestEditors: []common.RequestEditorFn{
			common.NewSessionHeaderRequestEditor("icSessionId"),
		},
	}

	apiClient := NewClientWithResponses(config, editors)
	return IdmcAdminV2Api{
		Client: apiClient,
	}

}
