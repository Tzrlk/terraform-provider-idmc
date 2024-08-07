package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/provider/models"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
	. "terraform-provider-idmc/internal/provider/utils"
	. "terraform-provider-idmc/internal/utils"
)

var _ ResourceWithConfigure = &RoleResource{}

type RoleResource struct {
	*IdmcProviderResource
}

func NewRoleResource() Resource {
	return &RoleResource{
		&IdmcProviderResource{},
	}
}

type RoleResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Privileges         types.List   `tfsdk:"privileges"`
	OrgId              types.String `tfsdk:"org_id"`
	DisplayName        types.String `tfsdk:"display_name"`
	DisplayDescription types.String `tfsdk:"display_description"`
	SystemRole         types.Bool   `tfsdk:"system_role"`
	Status             types.String `tfsdk:"status"`
	CreatedBy          types.String `tfsdk:"created_by"`
	UpdatedBy          types.String `tfsdk:"updated_by"`
	CreatedTime        types.String `tfsdk:"created_time"`
	UpdatedTime        types.String `tfsdk:"updated_time"`
}

// Metadata <editor-fold desc="Metadata" defaultstate="collapsed">
func (r RoleResource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// </editor-fold>

// Schema <editor-fold desc="Schema" defaultstate="collapsed">
func (r RoleResource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform_rest_api_version_3_resources/roles.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Service generated identifier for the role.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the role.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the role.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"privileges": schema.ListNestedAttribute{
				Description: "The privileges assigned to the role.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Privilege ID.",
							Optional:    true,
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
			"org_id": schema.StringAttribute{
				Description: "ID of the organization the role belongs to.",
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
		},
	}
}

// </editor-fold>

// Create <editor-fold desc="Create" defaultstate="collapsed">
func (r RoleResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadCreate)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV3(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data models.RoleValue
	if diags.Append(req.Plan.Get(ctx, &data)) {
		return
	}

	// Convert privilege set

	rolePrivileges := models.NewRolePrivilegeListValueFromList(diags.AtName("privileges"), data.Privileges)
	if diags.HasError() {
		return
	}

	apiRes, apiErr := client.CreateRoleWithResponse(ctx, v3.CreateRoleJSONRequestBody{
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Privileges:  Ptr(rolePrivileges.GetIds(diags).ToSlice()),
	})
	if diags.HandleError(apiErr) {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode() != 201 {
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
			diags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 201))
		}
		return
	}

	respData := *apiRes.JSON201

	// Update the configured state so instabilities can be detected.
	data.Id = types.StringPointerValue(respData.Id)
	data.Name = types.StringPointerValue(respData.RoleName)
	data.Description = types.StringPointerValue(respData.Description)

	// Update derived values
	data.OrgId = types.StringPointerValue(respData.OrgId)
	data.DisplayName = types.StringPointerValue(respData.DisplayName)
	data.DisplayDescription = types.StringPointerValue(respData.DisplayDescription)
	data.SystemRole = types.BoolPointerValue(respData.SystemRole)
	data.Status = types.StringPointerValue((*string)(respData.Status))
	data.CreatedBy = types.StringPointerValue(respData.CreatedBy)
	data.UpdatedBy = types.StringPointerValue(respData.UpdatedBy)
	data.CreatedTime = types.StringPointerValue(respData.CreateTime)
	data.UpdatedTime = types.StringPointerValue(respData.UpdateTime)

	// NOTE: Create does not return any privileges in the response.

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data))

}

// </editor-fold>

// Create <editor-fold desc="Read" defaultstate="collapsed">
func (r RoleResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadRead)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV3(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RoleResourceModel
	diags.Append(req.State.Get(ctx, &data))
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{
		Expand: Ptr(v3.GetRolesParamsExpandPrivileges),
	}
	if !data.Id.IsNull() {
		params.Q = Ptr(fmt.Sprintf("roleId==\"%s\"", data.Id.ValueString()))
	} else if !data.Name.IsNull() {
		params.Q = Ptr(fmt.Sprintf("roleName==\"%s\"", data.Name.ValueString()))
		diags.AtName("id").WithTitle("Issue reading resource").AddWarning(
			"No id for the role found in state. Falling back to name: %s", data.Name.ValueString())
	} else {
		diags.AtName("id").AddError(
			"No id or name for the role found in state.")
		return
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRolesWithResponse(ctx, params)
	if diags.HandleError(apiErr) {
		return
	}

	if diags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 200)) {
		return
	}

	apiItems := *apiRes.JSON200
	if len(apiItems) == 0 {
		// No matching resources, so junk it.
		resp.State.RemoveResource(ctx)
		return
	} else if len(apiItems) != 1 {
		diags.AddError(
			"Only one item was expected in the api response, not %d",
			len(apiItems),
		)
		return
	}

	// Update the configured state so instabilities can be detected.
	data.Id = types.StringPointerValue(apiItems[0].Id)
	data.Name = types.StringPointerValue(apiItems[0].RoleName)
	data.Description = types.StringPointerValue(apiItems[0].Description)

	// Update derived values
	data.OrgId = types.StringPointerValue(apiItems[0].OrgId)
	data.DisplayName = types.StringPointerValue(apiItems[0].DisplayName)
	data.DisplayDescription = types.StringPointerValue(apiItems[0].DisplayDescription)
	data.SystemRole = types.BoolPointerValue(apiItems[0].SystemRole)
	data.Status = types.StringPointerValue((*string)(apiItems[0].Status))
	data.CreatedBy = types.StringPointerValue(apiItems[0].CreatedBy)
	data.UpdatedBy = types.StringPointerValue(apiItems[0].UpdatedBy)
	data.CreatedTime = types.StringPointerValue(apiItems[0].CreateTime)
	data.UpdatedTime = types.StringPointerValue(apiItems[0].UpdateTime)

	// Handle more sketchy data
	if data.setPrivileges(diags, apiItems[0].Privileges) {
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data))

}

// </editor-fold>

// Update <editor-fold desc="Update" defaultstate="collapsed">
func (r RoleResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadDelete)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV3(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan and extract privileges.
	var plan RoleResourceModel
	diags.Append(req.Plan.Get(ctx, &plan))

	// Load config from state for comparison.
	var state RoleResourceModel
	diags.Append(req.State.Get(ctx, &state))

	// Only check for errors here so we can see if there are any issues with
	// either data structure before breaking.
	if diags.HasError() {
		return
	}

	// Load privileges into sets from both configs.
	planPrivileges := plan.getPrivileges(diags)
	statePrivileges := state.getPrivileges(diags)
	if diags.HasError() {
		return
	}

	privDiags := diags.AtName("privileges")

	// Add all the privileges that need to be added
	if privs := planPrivileges.Without(statePrivileges); privs.Size() > 0 {
		apiRes, apiErr := client.AddRolePrivilegesWithResponse(
			ctx,
			plan.Id.ValueString(),
			&v3.AddRolePrivilegesParams{},
			v3.AddRolePrivilegesJSONRequestBody{
				Privileges: privs.ToSlice(),
			},
		)
		if privDiags.HandleError(apiErr) {
			return
		}

		// Handle error responses.
		if apiRes.StatusCode() != 200 {
			CheckApiErrorV3(privDiags,
				apiRes.JSON400,
				apiRes.JSON401,
				apiRes.JSON403,
				apiRes.JSON404,
				apiRes.JSON500,
				apiRes.JSON502,
				apiRes.JSON503,
			)
			if !privDiags.HasError() {
				privDiags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 200))
			}
			return
		}

	}

	// Remove all the privileges that need to be removed
	if privs := statePrivileges.Without(planPrivileges); privs.Size() > 0 {
		apiRes, apiErr := client.RemoveRolePrivilegesWithResponse(
			ctx,
			plan.Id.ValueString(),
			&v3.RemoveRolePrivilegesParams{},
			v3.RemoveRolePrivilegesJSONRequestBody{
				Privileges: privs.ToSlice(),
			},
		)
		if privDiags.HandleError(apiErr) {
			return
		}

		// Handle error responses.
		if apiRes.StatusCode() != 200 {
			CheckApiErrorV3(privDiags,
				apiRes.JSON400,
				apiRes.JSON401,
				apiRes.JSON403,
				apiRes.JSON404,
				apiRes.JSON500,
				apiRes.JSON502,
				apiRes.JSON503,
			)
			if !privDiags.HasError() {
				privDiags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 200))
			}
			return
		}

	}

	// Save update result back to state.
	diags.Append(resp.State.Set(ctx, &plan))

}

// </editor-fold>

// Delete <editor-fold desc="Delete" defaultstate="collapsed">
func (r RoleResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadDelete)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV3(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RoleResourceModel
	if diags.Append(req.State.Get(ctx, &data)) {
		return
	}

	apiRes, apiErr := client.DeleteRoleWithResponse(ctx, data.Id.ValueString(), &v3.DeleteRoleParams{})
	if diags.HandleError(apiErr) {
		return
	}

	if diags.HandleError(RequireHttpStatus(&apiRes.ClientResponse, 200, 204)) {
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data))

}

// </editor-fold>

func (r RoleResourceModel) getPrivileges(diags DiagsHandler) models.RolePrivilegeListValue {
	diags = diags.AtName("privileges")
	var result models.RolePrivilegeListValue
	return result.NewRolePrivilegeListValueFromList(diags, r.Privileges)
}

func (r RoleResourceModel) setPrivileges(diags DiagsHandler, items *[]v3.RolePrivilegeItem) bool {
	diags = diags.AtName("privileges")

	if items == nil {
		diags.WithTitle("Issue handling resource API response").AddWarning(
			"Expected role privilege data, but received nothing.")
		r.Privileges = types.SetNull(types.StringType)
		return diags.HasError()
	}

	privAttrs := make([]attr.Value, len(*items))
	for index, item := range *items {
		privAttrs[index] = types.StringPointerValue(item.Id)
	}

	privAttr := diags.SetValue(types.StringType, privAttrs)
	if diags.HasError() {
		r.Privileges = types.SetUnknown(types.StringType)
		return true
	}

	r.Privileges = privAttr
	return diags.HasError()

}

//func (r RoleResource) updateRoleState(
//	diags *diag.Diagnostics,
//	state *RoleResourceModel,
//	data
//)
