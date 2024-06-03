package idmc

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go/types"
)

type ApiV2AgentInstallerInfoResponse struct {
	Type                string `json:"@type"`
	DownloadUrl         string `json:"download_url"`
	InstallToken        string `json:"install_token"`
	ChecksumDownloadUrl string `json:"checksum_download_url"`
}

func (api *ApiV2) DoAgentInstallerInfo(
	diag     *diag.Diagnostics,
	platform string,
) *ApiV2AgentInstallerInfoResponse {

	apiItem := ApiItem[types.Nil, ApiV2AgentInstallerInfoResponse]{
		Api:    api.Root,
		Method: "GET",
		Path:   "api/v2/agent/installerInfo/%s",
	}

	response := apiItem.Call(diag, []any{platform}, nil, nil)
	if diag.HasError() {
		return nil
	}

	return response

}
