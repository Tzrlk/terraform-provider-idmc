package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc"
	"terraform-provider-idmc/internal/idmc/admin/v3"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &RbacRolesDataSource{}

func NewRbacRolesDataSource() datasource.DataSource {
	return &RbacRolesDataSource{}
}

type RbacRolesDataSource struct {
	Client *v3.ClientWithResponses
}

type RbacRolesDataSourceModel struct {
	RoleId           types.String `tfsdk:"role_id"`
	RoleName         types.String `tfsdk:"role_name"`
	ExpandPrivileges types.Bool   `tfsdk:"expand_privileges"`
	Roles            types.Map    `tfsdk:"roles"`
}

func (d *RbacRolesDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_rbac_roles"
}

func (d *RbacRolesDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html",

		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
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
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"roles": schema.MapNestedAttribute{
				MarkdownDescription: "The query results",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						//"id": schema.StringAttribute{
						//	MarkdownDescription: "",
						//	Computed:            true,
						//},
						"org_id": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"created_by": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"updated_by": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"create_time": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"update_time": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"role_name": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"display_description": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"system_role": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"privileges": schema.MapNestedAttribute{
							MarkdownDescription: "",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									//"id": schema.StringAttribute{
									//	MarkdownDescription: "",
									//	Computed:            true,
									//},
									"name": schema.StringAttribute{
										MarkdownDescription: "",
										Computed:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "",
										Computed:            true,
									},
									"service": schema.StringAttribute{
										MarkdownDescription: "",
										Computed:            true,
									},
									"status": schema.StringAttribute{
										MarkdownDescription: "",
										Computed:            true,
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

func (d *RbacRolesDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {

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
	d.Client = api.Admin.V3.Client

}

func (d *RbacRolesDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	diags := &resp.Diagnostics

	// Load the previous state if present.
	var config RbacRolesDataSourceModel
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{}
	if !config.RoleId.IsNull() {
		query := fmt.Sprintf("roleName==\"%s\"", config.RoleId.ValueString())
		params.Q = &query
	} else if !config.RoleName.IsNull() {
		query := fmt.Sprintf("roleId==\"%s\"", config.RoleName.ValueString())
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
