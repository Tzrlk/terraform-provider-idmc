package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"

	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSourceWithConfigure = &PrivilegeListDataSource{}

type PrivilegeListDataSource struct {
	*IdmcProviderDataSource
}

func NewPrivilegeListDataSource() DataSource {
	return &PrivilegeListDataSource{
		&IdmcProviderDataSource{},
	}
}

type PrivilegeListDataSourceModel struct {
	Status     types.String `tfsdk:"status"`
	Privileges types.List   `tfsdk:"privileges"`
}
type PrivilegeListDataSourceModelPrivilege struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Service     types.String `tfsdk:"service"`
	Status      types.String `tfsdk:"status"`
}

func (d *PrivilegeListDataSource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privilege_list"
}

func (d *PrivilegeListDataSource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/b2b-gateway/current-version/rest-api-reference/platform_rest_api_version_3_resources/privileges.html",
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Description: "Filters the results by status. Use 'All' to get more than enabled and default results.",
				Optional:    true,
			},
			"privileges": schema.ListNestedAttribute{
				Description: "The results of the privilege list request.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Privilege ID.",
							Optional:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the privilege.",
							Optional:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the privilege.",
							Computed:    true,
						},
						"service": schema.StringAttribute{
							Description: "Service the privilege applies to.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the privilege:\n\tEnabled: License to use the privilege is valid.\n\tDisabled: License to use the privilege has expired.\n\tUnassigned: No license to use this privilege.\n\tDefault: Privilege included by default.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

var privilegeDataItemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"service":     types.StringType,
		"status":      types.StringType,
	},
}

func (d *PrivilegeListDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgDataSourceBadRead)

	client := d.GetApiClientV3(diags, MsgDataSourceBadRead)
	if diags.HasError() {
		return
	}

	// Load the previous state if present.
	var config PrivilegeListDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.ListPrivilegesParams{}
	if !config.Status.IsNull() {
		params.Q = utils.Ptr(fmt.Sprintf("status==\"%s\"", config.Status.ValueString()))
	}

	// Perform the API request.
	apiRes, apiErr := client.ListPrivilegesWithResponse(ctx, params)
	if errHandler(apiErr); diags.HasError() {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode() != 200 {
		CheckApiErrorV3(diags,
			apiRes.JSON400,
			apiRes.JSON401,
			apiRes.JSON403,
			apiRes.JSON404,
			apiRes.JSON500,
			apiRes.JSON502,
			apiRes.JSON503,
		)
		if !diags.HasError() {
			errHandler(RequireHttpStatus(200, apiRes))
		}
		return
	}

	config.Privileges = convertPrivilegeListResponse(diags, path.Root("privileges"), apiRes.JSON200)

	// Update the state and add the result
	diags.Append(resp.State.Set(ctx, &config)...)

}

func convertPrivilegeListResponse(diags *diag.Diagnostics, path path.Path, items *[]v3.RolePrivilegeItem) types.List {
	if items == nil {
		return types.ListNull(privilegeDataItemType)
	}

	privileges := make([]attr.Value, len(*items))
	for index, item := range *items {
		itemPath := path.AtListIndex(index)
		privileges[index] = UnwrapObjectValue(diags, itemPath, privilegeDataItemType.AttrTypes, map[string]attr.Value{
			"id":          types.StringPointerValue(item.Id),
			"name":        types.StringPointerValue(item.Name),
			"description": types.StringPointerValue(item.Description),
			"service":     types.StringPointerValue(item.Service),
			"status":      types.StringPointerValue((*string)(item.Status)),
		})
	}

	return UnwrapListValue(diags, path, rolesDataRolesPrivilegeType, privileges)
}
