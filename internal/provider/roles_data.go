package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"

	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSourceWithConfigure = &RolesDataSource{}
var _ DataSourceWithConfigValidators = &RolesDataSource{}

type RolesDataSource struct {
	IdmcProviderDataSource
}

func NewRolesDataSource() DataSource {
	return &RolesDataSource{
		IdmcProviderDataSource{},
	}
}

type RolesDataSourceModel struct {
	Filter types.Object `tfsdk:"filter"`
	Roles  types.Map    `tfsdk:"roles"`
}
type RolesDataSourceModelFilter struct {
	RoleId           types.String `tfsdk:"role_id"`
	RoleName         types.String `tfsdk:"role_name"`
	ExpandPrivileges types.Bool   `tfsdk:"expand_privileges"`
}

func (d *RolesDataSource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

func (d *RolesDataSource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html",

		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Description: "Allows for results to be narrowed.",
				Attributes: map[string]schema.Attribute{
					"role_id": schema.StringAttribute{
						Description: "Unique identifier for the role.",
						Optional:    true,
					},
					"role_name": schema.StringAttribute{
						Description: "Name of the role.",
						Optional:    true,
					},
					"expand_privileges": schema.BoolAttribute{
						Description: "Returns the privileges associated with the role specified in the query filter.",
						Optional:    true,
					},
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"roles": schema.MapNestedAttribute{
				Description: "The query results",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						//"id": schema.StringAttribute{
						//	Description: "Role ID.",
						//	Computed:            true,
						//},
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
						"privileges": schema.MapNestedAttribute{
							Description: "The privileges assigned to the role.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									//"id": schema.StringAttribute{
									//	Description: "Privilege ID.",
									//	Computed:            true,
									//},
									"name": schema.StringAttribute{
										Description: "Name of the privilege.",
										Computed:    true,
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
										Description: "Status of the privilege (Enabled/Disabled).",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *RolesDataSource) ConfigValidators(ctx context.Context) []ConfigValidator {
	return []ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("filter").AtName("role_id"),
			path.MatchRoot("filter").AtName("role_name"),
		),
	}
}

var rolesDataRolesPrivilegeType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"service":     types.StringType,
		"status":      types.StringType,
	},
}
var rolesDataRolesType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
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
		"privileges": types.MapType{
			ElemType: rolesDataRolesPrivilegeType,
		},
	},
}

// FIXME: Can only expand privileges on single result requests, not on full set. Full set doesn't return privileges at all. Maybe split into two data sources?

func (d *RolesDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgDataSourceBadRead)

	client := d.GetApiClientV3(diags)
	if diags.HasError() {
		return
	}

	// Load the previous state if present.
	var config RolesDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{}
	if !config.Filter.IsNull() {
		filter := UnwrapDiag(diags, path.Root("filter"), func() (RolesDataSourceModelFilter, diag.Diagnostics) {
			result := RolesDataSourceModelFilter{}
			resultDiag := config.Filter.As(ctx, &result, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})
			return result, resultDiag
		})
		if !filter.RoleId.IsNull() {
			params.Q = utils.Ptr(fmt.Sprintf("roleId==\"%s\"", filter.RoleId.ValueString()))
		} else if !filter.RoleName.IsNull() {
			params.Q = utils.Ptr(fmt.Sprintf("roleName==\"%s\"", filter.RoleName.ValueString()))
		}
		if filter.ExpandPrivileges.ValueBool() {
			params.Expand = utils.Ptr("privileges")
		}
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRolesWithResponse(ctx, params)
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

	// Handle empty JSON responses
	if apiRes.JSON200 == nil {
		diags.AddError(
			"IDMC API bad response payload",
			"Expected JSON payload, got nil.",
		)
		return
	}
	resBody := *apiRes.JSON200

	rolesPath := path.Root("roles")

	// Convert response data into terraform types.
	newRoles := make(map[string]attr.Value)
	for _, responseRole := range resBody {
		roleId := *responseRole.Id
		rolePath := rolesPath.AtMapKey(roleId)
		rolePrivilegesPath := rolePath.AtName("privileges")

		newRolePrivileges := make(map[string]attr.Value)
		if responseRole.Privileges != nil {
			for _, responseRolePrivilege := range *responseRole.Privileges {
				rolePrivilegePath := rolePrivilegesPath.AtMapKey(*responseRolePrivilege.Id)
				newRolePrivileges[*responseRolePrivilege.Id] = UnwrapObjectValue(diags, rolePrivilegePath, rolesDataRolesPrivilegeType.AttrTypes, map[string]attr.Value{
					"name":        types.StringPointerValue(responseRolePrivilege.Name),
					"description": types.StringPointerValue(responseRolePrivilege.Description),
					"service":     types.StringPointerValue(responseRolePrivilege.Service),
					"status":      types.StringPointerValue(responseRolePrivilege.Status),
				})
			}
		}

		newRoles[roleId] = UnwrapObjectValue(diags, rolePath, rolesDataRolesType.AttrTypes, map[string]attr.Value{
			"name":                types.StringPointerValue(responseRole.RoleName),
			"display_name":        types.StringPointerValue(responseRole.DisplayName),
			"org_id":              types.StringPointerValue(responseRole.OrgId),
			"description":         types.StringPointerValue(responseRole.Description),
			"display_description": types.StringPointerValue(responseRole.DisplayDescription),
			"system_role":         types.BoolPointerValue(responseRole.SystemRole),
			"status":              types.StringPointerValue(responseRole.Status),
			"created_by":          types.StringPointerValue(responseRole.CreatedBy),
			"updated_by":          types.StringPointerValue(responseRole.UpdatedBy),
			"created_time":        UnwrapNewRFC3339PointerValue(diags, rolePath.AtName("created_time"), responseRole.CreateTime),
			"updated_time":        UnwrapNewRFC3339PointerValue(diags, rolePath.AtName("updated_time"), responseRole.UpdateTime),
			"privileges":          UnwrapMapValue(diags, rolePath.AtName("privileges"), rolesDataRolesPrivilegeType, newRolePrivileges),
		})

	}

	config.Roles = UnwrapMapValue(diags, rolesPath, rolesDataRolesType, newRoles)

	// Update the state and add the result
	diags.Append(resp.State.Set(ctx, &config)...)

}
