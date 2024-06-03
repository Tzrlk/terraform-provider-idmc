package provider

import (
	"context"
	"fmt"
	"terraform-provider-idmc/internal/idmc"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &V2AgentInstallerInfoDataSource{}

func NewV2AgentInstallerInfoDataSource() datasource.DataSource {
	return &V2AgentInstallerInfoDataSource{}
}

type V2AgentInstallerInfoDataSource struct {
	Api *idmc.ApiV2
}

type V2AgentInstallerInfoDataSourceModel struct {
	Platform            types.String `tfsdk:"platform"`
	DownloadUrl         types.String `tfsdk:"download_url"`
	InstallToken        types.String `tfsdk:"install_token"`
	ChecksumDownloadUrl types.String `tfsdk:"checksum_download_url"`
}

func (d *V2AgentInstallerInfoDataSource) Metadata(
		ctx context.Context,
		req datasource.MetadataRequest,
		resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v2_agent_installer_info"
}

func (d *V2AgentInstallerInfoDataSource) Schema(
		ctx context.Context,
		req datasource.SchemaRequest,
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

func (d *V2AgentInstallerInfoDataSource) Configure(
		ctx context.Context,
		req datasource.ConfigureRequest,
		resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(*idmc.Api)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Api = api.V2

}

func (d *V2AgentInstallerInfoDataSource) Read(
		ctx context.Context,
		req datasource.ReadRequest,
		resp *datasource.ReadResponse) {
	var config V2AgentInstallerInfoDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	info := d.Api.DoAgentInstallerInfo(&resp.Diagnostics, config.Platform.ValueString())

	config.DownloadUrl         = types.StringValue(info.DownloadUrl)
	config.InstallToken        = types.StringValue(info.InstallToken)
	config.ChecksumDownloadUrl = types.StringValue(info.ChecksumDownloadUrl)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}
