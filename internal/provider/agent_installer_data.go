package provider

import (
	"context"
	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSourceWithConfigure = &AgentInstallerDataSource{}

type AgentInstallerDataSource struct {
	*IdmcProviderDataSource
}

func NewAgentInstallerDataSource() DataSource {
	return &AgentInstallerDataSource{
		&IdmcProviderDataSource{},
	}
}

type AgentInstallerDataSourceModel struct {
	Platform            types.String `tfsdk:"platform"`
	DownloadUrl         types.String `tfsdk:"download_url"`
	InstallToken        types.String `tfsdk:"install_token"`
	ChecksumDownloadUrl types.String `tfsdk:"checksum_download_url"`
}

func (d *AgentInstallerDataSource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_installer"
}

func (d *AgentInstallerDataSource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-2-resources/agent.html",
		Attributes: map[string]schema.Attribute{
			"platform": schema.StringAttribute{
				Description: "Platform of the Secure Agent machine. Must be one of the following values:\nwin64\nlinux64",
				Optional:    true,
			},
			"download_url": schema.StringAttribute{
				Description: "The URL of the latest Secure Agent installer package.",
				Computed:    true,
			},
			"install_token": schema.StringAttribute{
				Description: "Token needed to install and register a Secure Agent.",
				Computed:    true,
			},
			"checksum_download_url": schema.StringAttribute{
				Description: "The URL of the CRC-32 SHA256 package checksum.",
				Computed:    true,
			},
		},
	}
}

func (d *AgentInstallerDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgDataSourceBadRead)

	client := d.GetApiClientV2(diags, MsgDataSourceBadRead)
	if diags.HasError() {
		return
	}

	// Load the previous state if present.
	var config AgentInstallerDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Perform the API request.
	apiRes, apiErr := client.GetAgentInstallerInfoWithResponse(ctx, config.Platform.ValueString())
	if errHandler(apiErr); diags.HasError() {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode() != 200 {
		CheckApiErrorV2(diags,
			apiRes.JSON400,
			apiRes.JSON401,
			apiRes.JSON403,
			apiRes.JSON404,
			apiRes.JSON500,
			apiRes.JSON502,
			apiRes.JSON503,
		)
		if !diags.HasError() {
			errHandler(RequireHttpStatus(200, &apiRes.ClientResponse))
		}
		return
	}

	// Convert response data into terraform types.
	config.DownloadUrl = types.StringPointerValue(apiRes.JSON200.DownloadUrl)
	config.InstallToken = types.StringPointerValue(apiRes.JSON200.InstallToken)
	config.ChecksumDownloadUrl = types.StringPointerValue(apiRes.JSON200.ChecksumDownloadUrl)

	// Save result back to state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}
