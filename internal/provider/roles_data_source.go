package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"terraform-provider-idmc/internal/idmc/admin/v3"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &RolesDataSource{}

func NewRolesDataSource() datasource.DataSource {
	return &RolesDataSource{}
}

type RolesDataSource struct {
	Client *v3.ClientWithResponses
}

type RolesDataSourceModel struct {
	RoleId           types.String `tfsdk:"role_id"`
	RoleName         types.String `tfsdk:"role_name"`
	ExpandPrivileges types.Bool   `tfsdk:"expand_privileges"`
	Results          types.List   `tfsdk:"results"`
}

func (d *RolesDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

func (d *RolesDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html",

		Attributes: map[string]schema.Attribute{
			"role_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the role.",
				Optional:            true,
			},
			"role_name": schema.StringAttribute{
				MarkdownDescription: "Name of the role.",
				Optional:            true,
			},
			"expand_privileges": schema.BoolAttribute{
				MarkdownDescription: "Returns the privileges associated with the role specified in the query filter.",
				Optional:            true,
			},
			"results": schema.ListAttribute{
				MarkdownDescription: "The query results",
				Computed:            true,
			},
		},
	}
}

func (d *RolesDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*v3.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Client = client

}

func (d *RolesDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse) {
	var config RolesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &v3.GetRolesParams{}
	if !config.RoleId.IsNull() {
		query := fmt.Sprintf("roleName==\"%s\"", config.RoleId.ValueString())
		params.Q = &query
	} else if !config.RoleName.IsNull() {
		query := fmt.Sprintf("roleId==\"%s\"", config.RoleName.ValueString())
		params.Q = &query
	}
	if config.ExpandPrivileges.ValueBool() {
		expand := v3.Privileges
		params.Expand = &expand
	}

	response, err := d.Client.GetRolesWithResponse(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Http Request Failure",
			fmt.Sprintf("IDMC Api request failure: %s", err),
		)
	}

	// TODO: response data.
	//config.Results = types.ListValue(types.String, &attr.Value[
	//	response.JSON200
	//])

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}
