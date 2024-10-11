package v3

import (
	"terraform-provider-idmc/internal/idmc/common"
)

//go:generate -command oapi-codegen go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
//go:generate oapi-codegen -config ./codegen.yml -generate types -o ./types.gen.go ./openapi.yml
//go:generate oapi-codegen -config ./codegen.yml -generate client -o ./client.gen.go ./openapi.yml

type IdmcAdminV3Api struct {
	Client         ClientWithResponses
	rolePrivileges *map[string]RolePrivilegeItem
}

func NewIdmcAdminV3Api(config *common.ClientConfig) IdmcAdminV3Api {

	editors := common.ClientConfigEditor{
		RequestEditors: []common.RequestEditorFn{
			common.WithSessionHeader("INFA-SESSION-ID"),
			common.WithRequestHeader("Accept", "application/json"),
		},
	}

	apiClient := NewClientWithResponses(config, editors)
	return IdmcAdminV3Api{
		Client: apiClient,
	}

}
