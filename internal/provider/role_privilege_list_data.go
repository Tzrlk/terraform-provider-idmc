package provider

import (
	"context"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	. "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
	"terraform-provider-idmc/internal/idmc/v3"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ DataSourceWithConfigure = &RolePrivilegeListDataSource{}

type RolePrivilegeListDataSource struct {
	*IdmcProviderDataSource
}

func NewRolePrivilegeListDataSource() DataSource {
	return &RolePrivilegeListDataSource{
		&IdmcProviderDataSource{},
	}
}

type RolePrivilegeListDataSourceModel struct {
	Status     types.String `tfsdk:"status"`
	Privileges types.List   `tfsdk:"privileges"`
}
type RolePrivilegeListDataSourceModelPrivilege struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Service     types.String `tfsdk:"service"`
	Status      types.String `tfsdk:"status"`
}

func (d *RolePrivilegeListDataSource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_privilege_list"
}

func (d *RolePrivilegeListDataSource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
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

func (d *RolePrivilegeListDataSource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgDataSourceBadRead)
	defer func() { diags.HandlePanic(recover()) }()

	apiV3 := d.GetApi().V3

	// Load the previous state if present.
	var config RolePrivilegeListDataSourceModel
	if diags.Append(req.Config.Get(ctx, &config)) {
		return
	}

	// Load up default query privileges
	privMap, err := apiV3.GetRolePrivileges(ctx, config.Status.ValueStringPointer())
	if diags.HandleError(err) {
		return
	}

	// Set the privileges from the map values.
	if config.setPrivileges(diags, lo.ToPtr(maps.Values(privMap))) {
		return
	}

	// Update the state and add the result
	diags.Append(resp.State.Set(ctx, &config))

}

func (r *RolePrivilegeListDataSourceModel) setPrivileges(diags DiagsHandler, items *[]v3.RolePrivilegeItem) bool {
	diags = diags.AtName("privileges")

	if items == nil {
		diags.AddWarning("Expected API response to contain privilege list.")
		r.Privileges = types.ListNull(privilegeDataItemType)
		return false
	}

	privileges := make([]attr.Value, len(*items))
	for index, item := range *items {
		privileges[index] = diags.AtListIndex(index).ObjectValue(privilegeDataItemType.AttrTypes, map[string]attr.Value{
			"id":          types.StringValue(item.Id),
			"name":        types.StringValue(item.Name),
			"description": types.StringPointerValue(item.Description),
			"service":     types.StringValue(item.Service),
			"status":      types.StringValue((string)(item.Status)),
		})
	}

	privAttr := diags.ListValue(rolesDataRolesPrivilegeType, privileges)
	if diags.HasError() {
		return true
	}

	r.Privileges = privAttr
	return false

}
