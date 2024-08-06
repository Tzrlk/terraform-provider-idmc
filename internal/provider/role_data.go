package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"

	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSourceWithConfigure = &RoleDataSource{}
var _ DataSourceWithConfigValidators = &RoleDataSource{}

type RoleDataSource struct {
	*IdmcProviderDataSource
}

func NewRoleDataSource() DataSource {
	return &RoleDataSource{
		&IdmcProviderDataSource{},
	}
}

type RoleDataSourceModel struct {
	Id                 types.String      `tfsdk:"id"`
	Name               types.String      `tfsdk:"name"`
	OrgId              types.String      `tfsdk:"org_id"`
	DisplayName        types.String      `tfsdk:"display_name"`
	Description        types.String      `tfsdk:"description"`
	DisplayDescription types.String      `tfsdk:"display_description"`
	SystemRole         types.Bool        `tfsdk:"system_role"`
	Status             types.String      `tfsdk:"status"`
	CreatedBy          types.String      `tfsdk:"created_by"`
	UpdatedBy          types.String      `tfsdk:"updated_by"`
	CreatedTime        timetypes.RFC3339 `tfsdk:"created_time"`
	UpdatedTime        timetypes.RFC3339 `tfsdk:"updated_time"`
	Privileges         types.List        `tfsdk:"privileges"`
}
type RoleDataSourceModelPrivilege struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Service     types.String `tfsdk:"service"`
	Status      types.String `tfsdk:"status"`
}

func (d *RoleDataSource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *RoleDataSource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Role ID.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the role.",
				Optional:    true,
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
			"privileges": schema.ListNestedAttribute{
				Description: "The privileges assigned to the role.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Privilege ID.",
							Computed:    true,
						},
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
	}
}

func (d *RoleDataSource) ConfigValidators(_ context.Context) []ConfigValidator {
	return []ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

var rolesDataRolesPrivilegeType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"service":     types.StringType,
		"status":      types.StringType,
	},
}

func (d *RoleDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := NewDiagsHandler(&resp.Diagnostics, MsgDataSourceBadRead)
	defer func() { diags.HandlePanic(recover()) }()

	client := d.GetApiClientV3(diags)
	if diags.HasError() {
		return
	}

	// Load the previous state if present.
	var config RoleDataSourceModel
	if diags.Append(req.Config.Get(ctx, &config)) {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{
		Expand: utils.Ptr(v3.GetRolesParamsExpandPrivileges),
	}
	if !config.Id.IsNull() {
		params.Q = utils.Ptr(fmt.Sprintf("roleId==\"%s\"", config.Id.ValueString()))
	} else if !config.Name.IsNull() {
		params.Q = utils.Ptr(fmt.Sprintf("roleName==\"%s\"", config.Name.ValueString()))
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRolesWithResponse(ctx, params)
	if diags.HandleError(apiErr) {
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
			diags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 200))
		}
		return
	}

	if apiRes.JSON200 == nil || len(*apiRes.JSON200) < 1 {
		if !config.Id.IsNull() {
			diags.AddError("no results returned for given id: %s", config.Id.ValueString())
		} else if !config.Name.IsNull() {
			diags.AddError("no results returned for given name: %s", config.Name.ValueString())
		}
		return
	}

	item := (*apiRes.JSON200)[0]
	config.Id = types.StringPointerValue(item.Id)
	config.Name = types.StringPointerValue(item.RoleName)
	config.DisplayName = types.StringPointerValue(item.DisplayName)
	config.OrgId = types.StringPointerValue(item.OrgId)
	config.Description = types.StringPointerValue(item.Description)
	config.DisplayDescription = types.StringPointerValue(item.DisplayDescription)
	config.SystemRole = types.BoolPointerValue(item.SystemRole)
	config.Status = types.StringPointerValue((*string)(item.Status))
	config.CreatedBy = types.StringPointerValue(item.CreatedBy)
	config.UpdatedBy = types.StringPointerValue(item.UpdatedBy)
	config.CreatedTime = diags.AtName("created_time").TimePointer(item.CreateTime)
	config.UpdatedTime = diags.AtName("updated_time").TimePointer(item.UpdateTime)

	// Handle more sketchy config data
	if config.setPrivileges(diags, item.Privileges) {
		return
	}

	// Update the state and add the result
	diags.Append(resp.State.Set(ctx, &config))

}

func (r *RoleDataSourceModel) setPrivileges(diags DiagsHandler, items *[]v3.RolePrivilegeItem) bool {
	diags = diags.AtName("privileges")

	if items == nil {
		diags.WithTitle("Issue reading datasource.").AddWarning(
			"Expected API response to contain role list.")
		r.Privileges = types.ListNull(rolesDataRolesPrivilegeType)
		return false
	}

	privAttrs := make([]attr.Value, len(*items))
	for index, item := range *items {
		privAttrs[index] = diags.AtListIndex(index).ObjectValue(rolesDataRolesPrivilegeType.AttrTypes, map[string]attr.Value{
			"id":          types.StringPointerValue(item.Id),
			"name":        types.StringPointerValue(item.Name),
			"description": types.StringPointerValue(item.Description),
			"service":     types.StringPointerValue(item.Service),
			"status":      types.StringPointerValue((*string)(item.Status)),
		})
	}

	privAttr := diags.ListValue(rolesDataRolesPrivilegeType, privAttrs)
	if diags.HasError() {
		return true
	}

	r.Privileges = privAttr
	return false

}
