package datasources

import (
	"context"
	"fmt"
	"terraform-provider-idmc/internal/idmc"
	v2 "terraform-provider-idmc/internal/idmc/admin/v2"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AgentInstallerDataSource{}

func NewAgentInstallerDataSource() datasource.DataSource {
	return &AgentInstallerDataSource{}
}

type AgentInstallerDataSource struct {
	Client *v2.ClientWithResponses
}

type AgentInstallerDataSourceModel struct {
	Platform            types.String `tfsdk:"platform"`
	DownloadUrl         types.String `tfsdk:"download_url"`
	InstallToken        types.String `tfsdk:"install_token"`
	ChecksumDownloadUrl types.String `tfsdk:"checksum_download_url"`
}

func (d *AgentInstallerDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v2_agent_installer_info"
}

func (d *AgentInstallerDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform-rest-api-version-2-resources/agent.html",

		Attributes: map[string]schema.Attribute{
			"platform": schema.StringAttribute{
				MarkdownDescription: "Platform of the Secure Agent machine. Must be one of the following values:\nwin64\nlinux64",
				Optional:            true,
				// TODO: Implement validation.
			},
			"download_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the latest Secure Agent installer package.",
				Computed:            true,
			},
			"install_token": schema.StringAttribute{
				MarkdownDescription: "Token needed to install and register a Secure Agent.",
				Computed:            true,
			},
			"checksum_download_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the CRC-32 SHA256 package checksum.",
				Computed:            true,
			},
		},
	}
}

func (d *AgentInstallerDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(*idmc.IdmcApi)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Client = api.Admin.V2.Client

}

func (d *AgentInstallerDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	diags := &resp.Diagnostics

	// Load the previous state if present.
	var config AgentInstallerDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Perform the API request.
	response, err := d.Client.GetAgentInstallerInfoWithResponse(ctx, config.Platform.ValueString())
	if err != nil {
		diags.AddError(
			"Http Request Failure",
			fmt.Sprintf("IDMC Api request failure: %s", err),
		)
		return
	}

	// Convert response data into terraform types.
	config.DownloadUrl = types.StringPointerValue(response.JSON200.DownloadUrl)
	config.InstallToken = types.StringPointerValue(response.JSON200.InstallToken)
	config.ChecksumDownloadUrl = types.StringPointerValue(response.JSON200.ChecksumDownloadUrl)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}
