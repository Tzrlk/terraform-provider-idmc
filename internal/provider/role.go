package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
	"terraform-provider-idmc/internal/idmc/v3"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
	. "terraform-provider-idmc/internal/provider/utils"
	. "terraform-provider-idmc/internal/utils"
)

var _ ResourceWithConfigure = &RoleResource{}
var _ ResourceWithImportState = &RoleResource{}

type RoleResource struct {
	IdmcProviderResource
}

func NewRoleResource() Resource {
	return &RoleResource{
		IdmcProviderResource{
			Name: "role",
		},
	}
}

type RoleResourceModel struct {
	Id                 types.String      `tfsdk:"id"`
	Name               types.String      `tfsdk:"name"`
	Description        types.String      `tfsdk:"description"`
	Privileges         types.Set         `tfsdk:"privileges"`
	OrgId              types.String      `tfsdk:"org_id"`
	DisplayName        types.String      `tfsdk:"display_name"`
	DisplayDescription types.String      `tfsdk:"display_description"`
	SystemRole         types.Bool        `tfsdk:"system_role"`
	Status             types.String      `tfsdk:"status"`
	CreatedBy          types.String      `tfsdk:"created_by"`
	UpdatedBy          types.String      `tfsdk:"updated_by"`
	CreatedTime        timetypes.RFC3339 `tfsdk:"created_time"`
	UpdatedTime        timetypes.RFC3339 `tfsdk:"updated_time"`
}

// Schema <editor-fold desc="Schema" defaultstate="collapsed">
func (r *RoleResource) Schema(ctx context.Context, req SchemaRequest, resp *SchemaResponse) {
	r.IdmcProviderResource.Schema(ctx, req, resp)
	resp.Schema.Description = "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform_rest_api_version_3_resources/roles.html"
	resp.Schema.Attributes = MapMerge(resp.Schema.Attributes, map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "Name of the role.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"description": schema.StringAttribute{
			Description: "Description of the role.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"privileges": schema.SetAttribute{
			Description: "The privileges assigned to the role.",
			Required:    true,
			ElementType: types.StringType,
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
	})
}

// </editor-fold>

// Create <editor-fold desc="Create" defaultstate="collapsed">
func (r *RoleResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadCreate)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApi().V3.Client

	// Load configuration from plan.
	var data RoleResourceModel
	if diags.Append(req.Plan.Get(ctx, &data)) {
		return
	}

	// Convert privilege set
	rolePrivileges := data.getPrivileges(diags)
	if diags.HasError() {
		return
	}

	apiRes, apiErr := client.CreateRoleWithResponse(ctx, v3.CreateRoleJSONRequestBody{
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Privileges:  lo.ToPtr(rolePrivileges.ToSlice()),
	})
	if diags.HandleError(apiErr) {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode != 201 {
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
	data.CreatedTime = diags.TimePointer(respData.CreateTime)
	data.UpdatedTime = diags.TimePointer(respData.UpdateTime)

	// NOTE: Create does not return any privileges in the response.

	// If we had trouble parsing the data.
	if diags.HasError() {
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data))

}

// </editor-fold>

// Read <editor-fold desc="Read" defaultstate="collapsed">
func (r *RoleResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadRead)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApi().V3.Client

	// Load configuration from plan.
	var data RoleResourceModel
	diags.Append(req.State.Get(ctx, &data))
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{
		Expand: lo.ToPtr(v3.GetRolesParamsExpandPrivileges),
	}
	if !data.Id.IsNull() {
		params.Q = lo.ToPtr(fmt.Sprintf("roleId==\"%s\"", data.Id.ValueString()))
	} else if !data.Name.IsNull() {
		params.Q = lo.ToPtr(fmt.Sprintf("roleName==\"%s\"", data.Name.ValueString()))
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
	data.CreatedTime = diags.TimePointer(apiItems[0].CreateTime)
	data.UpdatedTime = diags.TimePointer(apiItems[0].UpdateTime)

	// Handle more sketchy data
	data.Privileges = diags.SetValuePointerFromFn(types.StringType, func() *[]attr.Value {
		if apiItems[0].Privileges == nil {
			diags.WithTitle("Issue handling resource API response").AddWarning(
				"Expected role privilege data, but received nothing.")
			return nil
		}

		return lo.ToPtr(lo.Map(*apiItems[0].Privileges, func(item v3.RolePrivilegeItem, index int) attr.Value {
			return types.StringValue(item.Id)
		}))
	})

	// If we had trouble parsing the data.
	if diags.HasError() {
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data))

}

// </editor-fold>

// Update <editor-fold desc="Update" defaultstate="collapsed">
func (r *RoleResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadDelete)
	defer func() { diags.HandlePanic(recover()) }()

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
	privsToAdd := planPrivileges.Without(statePrivileges)
	if privsToAdd.Size() > 0 {
		err := r.Api.V3.AddRolePrivileges(ctx, plan.Id.ValueString(), privsToAdd.ToSlice())
		if privDiags.HandleError(err) {
			return
		}
		// Save update result back to state.
		state.Privileges = diags.SetValueFromFn(types.StringType, func() []attr.Value {
			return lo.Map(statePrivileges.Union(privsToAdd).ToSlice(), func(from string, index int) attr.Value {
				return types.StringValue(from)
			})
		})
		if diags.HasError() {
			return
		}
		// Update the state so a failure to remove doesn't break things.
		diags.Append(resp.State.Set(ctx, &state))
	}

	// Remove all the privileges that need to be removed
	privsToRem := statePrivileges.Without(planPrivileges)
	if privsToRem.Size() > 0 {
		err := r.Api.V3.RemoveRolePrivileges(ctx, plan.Id.ValueString(), privsToRem.ToSlice())
		if privDiags.HandleError(err) {
			return
		}
		state.Privileges = plan.Privileges
		// Save update result back to state.
		diags.Append(resp.State.Set(ctx, &state))
	}

}

// </editor-fold>

// Delete <editor-fold desc="Delete" defaultstate="collapsed">
func (r *RoleResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := NewDiagsHandler(ctx, &resp.Diagnostics, MsgResourceBadDelete)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApi().V3.Client

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

func (r *RoleResourceModel) getPrivileges(diags DiagsHandler) *HashSet[string] {
	privilegesPath := path.Root("privileges")
	return NewHashSetAfter(func(set *HashSet[string]) {
		for _, element := range r.Privileges.Elements() {
			elementAttr, castOk := element.(types.String)
			if castOk && !elementAttr.IsNull() && !elementAttr.IsUnknown() {
				set.Add(elementAttr.ValueString())
				continue
			}
			diags.WithPath(privilegesPath.AtSetValue(element)).AddError(
				"Encountered a bad value loading set data: %s", element)
		}
	})
}

//func (r *RoleResource) updateRoleState(
//	diags *diag.Diagnostics,
//	state *RoleResourceModel,
//	data
//)
