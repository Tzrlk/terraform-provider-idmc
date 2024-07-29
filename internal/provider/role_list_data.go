package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v3"

	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSourceWithConfigure = &RoleListDataSource{}

type RoleListDataSource struct {
	*IdmcProviderDataSource
}

func NewRoleListDataSource() DataSource {
	return &RoleListDataSource{
		&IdmcProviderDataSource{},
	}
}

type RoleListDataSourceModel struct {
	Roles types.List `tfsdk:"roles"`
}

func (d *RoleListDataSource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_list"
}

func (d *RoleListDataSource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html",
		Attributes: map[string]schema.Attribute{
			"roles": schema.ListNestedAttribute{
				Description: "The query results",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Role ID.",
							Computed:    true,
						},
						"org_id": schema.StringAttribute{
							Description: "ID of the organization the role belongs to.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "User who created the role.",
							Computed:    true,
						},
						"updated_by": schema.StringAttribute{
							Description: "User who last updated the role.",
							Computed:    true,
						},
						"created_time": schema.StringAttribute{
							CustomType:  timetypes.RFC3339Type{},
							Description: "Date and time the role was created.",
							Computed:    true,
						},
						"updated_time": schema.StringAttribute{
							CustomType:  timetypes.RFC3339Type{},
							Description: "Date and time the role was last updated.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the role.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the role.",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Role name displayed in the user interface.",
							Computed:    true,
						},
						"display_description": schema.StringAttribute{
							Description: "Description displayed in the user interface.",
							Computed:    true,
						},
						"system_role": schema.BoolAttribute{
							Description: "Whether the role is a system-defined role.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Whether the organization's license to use the role is valid or has expired.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

var roleListDataRoleType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                  types.StringType,
		"name":                types.StringType,
		"display_name":        types.StringType,
		"org_id":              types.StringType,
		"description":         types.StringType,
		"display_description": types.StringType,
		"system_role":         types.BoolType,
		"status":              types.StringType,
		"created_by":          types.StringType,
		"updated_by":          types.StringType,
		"created_time":        timetypes.RFC3339Type{},
		"updated_time":        timetypes.RFC3339Type{},
	},
}

func (d *RoleListDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgDataSourceBadRead)

	client := d.GetApiClientV3(diags, MsgDataSourceBadRead)
	if diags.HasError() {
		return
	}

	// Load the previous state if present.
	var config RoleListDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRolesWithResponse(ctx, &v3.GetRolesParams{})
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

	config.Roles = convertRoleListResponse(diags, apiRes.JSON200)

	// Update the state and add the result
	diags.Append(resp.State.Set(ctx, &config)...)

}

func convertRoleListResponse(diags *diag.Diagnostics, items *[]v3.GetRolesResponseBodyItem) types.List {
	if items == nil {
		return types.ListUnknown(roleListDataRoleType)
	}

	roles := make([]attr.Value, len(*items))
	rolesPath := path.Root("roles")
	for index, item := range *items {
		rolePath := rolesPath.AtListIndex(index)
		roles[index] = UnwrapObjectValue(diags, rolePath, roleListDataRoleType.AttrTypes, map[string]attr.Value{
			"id":                  types.StringPointerValue(item.Id),
			"name":                types.StringPointerValue(item.RoleName),
			"display_name":        types.StringPointerValue(item.DisplayName),
			"org_id":              types.StringPointerValue(item.OrgId),
			"description":         types.StringPointerValue(item.Description),
			"display_description": types.StringPointerValue(item.DisplayDescription),
			"system_role":         types.BoolPointerValue(item.SystemRole),
			"status":              types.StringPointerValue((*string)(item.Status)),
			"created_by":          types.StringPointerValue(item.CreatedBy),
			"updated_by":          types.StringPointerValue(item.UpdatedBy),
			"created_time":        UnwrapNewRFC3339PointerValue(diags, rolePath.AtName("created_time"), item.CreateTime),
			"updated_time":        UnwrapNewRFC3339PointerValue(diags, rolePath.AtName("updated_time"), item.UpdateTime),
		})
	}

	return UnwrapListValue(diags, rolesPath, roleListDataRoleType, roles)
}
