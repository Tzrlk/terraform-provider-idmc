package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v2"

	. "github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ DataSource = &AgentInstallerDataSource{}

func NewAgentInstallerDataSource() DataSource {
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

func (d *AgentInstallerDataSource) Configure(_ context.Context, req ConfigureRequest, resp *ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*IdmcProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Client = data.Api.V2.Client

}

func (d *AgentInstallerDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
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
