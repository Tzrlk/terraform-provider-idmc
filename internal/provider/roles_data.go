package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/admin/v3"

	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSource = &RolesDataSource{}
var _ DataSourceWithConfigValidators = &RolesDataSource{}

func NewRolesDataSource() DataSource {
	return &RolesDataSource{}
}

type RolesDataSource struct {
	Client *v3.ClientWithResponses
}

type RolesDataSourceModel struct {
	RoleId           types.String `tfsdk:"role_id"`
	RoleName         types.String `tfsdk:"role_name"`
	ExpandPrivileges types.Bool   `tfsdk:"expand_privileges"`
	Roles            types.Map    `tfsdk:"roles"`
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
							Description: "Date and time the role was created.",
							Computed:    true,
						},
						"updated_time": schema.StringAttribute{
							Description: "Date and time the role was last updated.",
							Computed:    true,
						},
						"role_name": schema.StringAttribute{
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
						"system_role": schema.StringAttribute{
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

func (d *RolesDataSource) Configure(_ context.Context, req ConfigureRequest, resp *ConfigureResponse) {

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
	d.Client = data.Api.Admin.V3.Client

}

func (d *RolesDataSource) ConfigValidators(ctx context.Context) []ConfigValidator {
	return []ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("filter").AtName("role_id"),
			path.MatchRoot("filter").AtName("role_name"),
		),
	}
}

func (d *RolesDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics

	// Load the previous state if present.
	var config RolesDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{}
	if !config.RoleId.IsNull() {
		query := fmt.Sprintf("roleId==\"%s\"", config.RoleId.ValueString())
		params.Q = &query
	} else if !config.RoleName.IsNull() {
		query := fmt.Sprintf("roleName==\"%s\"", config.RoleName.ValueString())
		params.Q = &query
	}
	if config.ExpandPrivileges.ValueBool() {
		expand := "privileges"
		params.Expand = &expand
	}

	// Perform the API request.
	response, err := d.Client.GetRolesWithResponse(ctx, params)
	if err != nil {
		diags.AddError(
			"Http Request Failure",
			fmt.Sprintf("IDMC Api request failure: %s", err),
		)
		return
	}

	rolesPrivilegeType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":        types.StringType,
			"description": types.StringType,
			"service":     types.StringType,
			"status":      types.StringType,
		},
	}

	rolesPath := path.Root("roles")
	rolesType := types.ObjectType{
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
			"created_time":        types.StringType,
			"updated_time":        types.StringType,
			"privileges": types.MapType{
				ElemType: rolesPrivilegeType,
			},
		},
	}

	// Convert response data into terraform types.
	newRoles := make(map[string]attr.Value)
	for _, responseRole := range *response.JSON200 {
		rolePath := rolesPath.AtMapKey(*responseRole.Id)
		rolePrivilegesPath := rolePath.AtName("privileges")

		newRolePrivileges := make(map[string]attr.Value)
		for _, responseRolePrivilege := range *responseRole.Privileges {
			rolePrivilegePath := rolePrivilegesPath.AtMapKey(*responseRolePrivilege.Id)
			newRolePrivileges[*responseRolePrivilege.Id] = UnwrapObjectValue(diags, rolePrivilegePath, rolesPrivilegeType.AttrTypes, map[string]attr.Value{
				"name":        types.StringPointerValue(responseRolePrivilege.Name),
				"description": types.StringPointerValue(responseRolePrivilege.Description),
				"service":     types.StringPointerValue(responseRolePrivilege.Service),
				"status":      types.StringPointerValue(responseRolePrivilege.Status),
			})
		}

		newRoles[*responseRole.Id] = UnwrapObjectValue(diags, rolePath, rolesType.AttrTypes, map[string]attr.Value{
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
			"privileges":          UnwrapMapValue(diags, rolePath.AtName("privileges"), rolesPrivilegeType, newRolePrivileges),
		})

	}

	config.Roles = UnwrapMapValue(diags, rolesPath, rolesType, newRoles)

	// Update the state and add the result
	diags.Append(resp.State.Set(ctx, &config)...)

}
