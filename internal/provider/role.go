package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
	"terraform-provider-idmc/internal/idmc/v3"

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
	Privileges         types.Set    `tfsdk:"privileges"`
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

func (r RoleResource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r RoleResource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "TODO",
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

func (r RoleResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadCreate)

	client := r.GetApiClientV3(diags, MsgResourceBadCreate)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RoleResourceModel
	diags.Append(req.Plan.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	// Convert privilege set
	rolePrivilegeElements := data.Privileges.Elements()
	rolePrivilegeValues := make([]string, len(rolePrivilegeElements))
	for index, rolePrivilegeElement := range rolePrivilegeElements {
		rolePrivilegeValue, ok := rolePrivilegeElement.(types.String)
		if !ok {
			diags.AddAttributeError(
				path.Root("privileges").AtSetValue(rolePrivilegeElement),
				"Bad element type in privileges",
				fmt.Sprintf("An item in the privilege list isn't a string: %s", rolePrivilegeElement.String()),
			)
		}
		rolePrivilegeValues[index] = rolePrivilegeValue.String()
	}
	if diags.HasError() {
		return
	}

	apiRes, apiErr := client.CreateRoleWithResponse(ctx, v3.CreateRoleJSONRequestBody{
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Privileges:  &rolePrivilegeValues,
	})
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
			errHandler(RequireHttpStatus(200, &apiRes.ClientResponse))
		}
		return
	}

	// Update the configured state so instabilities can be detected.
	data.Id = types.StringPointerValue(apiRes.JSON200.Id)
	data.Name = types.StringPointerValue(apiRes.JSON200.RoleName)
	data.Description = types.StringPointerValue(apiRes.JSON200.Description)
	data.Privileges = UnwrapDiag(diags, path.Root("privileges"), func() (types.Set, diag.Diagnostics) {
		apiRolePrivilegeItems := *apiRes.JSON200.Privileges
		apiRolePrivilegeAttrs := make([]attr.Value, len(apiRolePrivilegeItems))
		for index, apiRolePrivilegeItem := range apiRolePrivilegeItems {
			apiRolePrivilegeAttrs[index] = types.StringPointerValue(apiRolePrivilegeItem.Id)
		}
		return types.SetValue(types.StringType, apiRolePrivilegeAttrs)
	})

	// Update derived values
	data.OrgId = types.StringPointerValue(apiRes.JSON200.OrgId)
	data.DisplayName = types.StringPointerValue(apiRes.JSON200.DisplayName)
	data.DisplayDescription = types.StringPointerValue(apiRes.JSON200.DisplayDescription)
	data.SystemRole = types.BoolPointerValue(apiRes.JSON200.SystemRole)
	data.Status = types.StringPointerValue((*string)(apiRes.JSON200.Status))
	data.CreatedBy = types.StringPointerValue(apiRes.JSON200.CreatedBy)
	data.UpdatedBy = types.StringPointerValue(apiRes.JSON200.UpdatedBy)
	data.CreatedTime = types.StringPointerValue(apiRes.JSON200.CreateTime)
	data.UpdatedTime = types.StringPointerValue(apiRes.JSON200.UpdateTime)

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

func (r RoleResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadRead)

	client := r.GetApiClientV3(diags, MsgResourceBadRead)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RoleResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	// Obtain request parameters from config.
	params := &v3.GetRolesParams{}
	if !data.Id.IsNull() {
		params.Q = Ptr(fmt.Sprintf("roleId==\"%s\"", data.Id.ValueString()))
	} else if !data.Name.IsNull() {
		params.Q = Ptr(fmt.Sprintf("roleName==\"%s\"", data.Name.ValueString()))
		diags.AddAttributeWarning(
			path.Root("id"),
			"Resource refresh required sketchy fall-back",
			fmt.Sprintf(
				"No id for the role found in state. Falling back to name: %s",
				data.Name.ValueString(),
			),
		)
	} else {
		diags.AddAttributeError(
			path.Root("id"),
			"Unable to Refresh Resource",
			"No id or name for the role found in state.",
		)
		return
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRolesWithResponse(ctx, params)
	if errHandler(apiErr); diags.HasError() {
		return
	}

	errHandler(RequireHttpStatus(200, &apiRes.ClientResponse))
	if diags.HasError() {
		return
	}

	apiItems := *apiRes.JSON200
	if len(apiItems) == 0 {
		// No matching resources, so junk it.
		resp.State.RemoveResource(ctx)
		return
	} else if len(apiItems) != 1 {
		diags.AddError(
			"Bad response from API",
			fmt.Sprintf(
				"Only one item was expected in the api response, not %d",
				len(apiItems),
			),
		)
		return
	}

	// Update the configured state so instabilities can be detected.
	data.Id = types.StringPointerValue(apiItems[0].Id)
	data.Name = types.StringPointerValue(apiItems[0].RoleName)
	data.Description = types.StringPointerValue(apiItems[0].Description)
	data.Privileges = UnwrapDiag(diags, path.Root("privileges"), func() (types.Set, diag.Diagnostics) {
		apiRolePrivilegeItems := *apiItems[0].Privileges
		apiRolePrivilegeAttrs := make([]attr.Value, len(apiRolePrivilegeItems))
		for index, apiRolePrivilegeItem := range apiRolePrivilegeItems {
			apiRolePrivilegeAttrs[index] = types.StringPointerValue(apiRolePrivilegeItem.Id)
		}
		return types.SetValue(types.StringType, apiRolePrivilegeAttrs)
	})

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

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

func (r RoleResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadUpdate)

	client := r.GetApiClientV3(diags, MsgResourceBadUpdate)
	if diags.HasError() {
		return
	}

	// Load config from state for comparison.
	var state RoleResourceModel
	diags.Append(req.State.Get(ctx, &state)...)

	// Load configuration from plan.
	var plan RoleResourceModel
	diags.Append(req.Plan.Get(ctx, &plan)...)

	// Only check for errors here so we can see if there are any issues with
	// either data structure before breaking.
	if diags.HasError() {
		return
	}

	rolePrivilegesPath := path.Root("privileges")

	// Load all planned privileges as if they're all to be added.
	privilegesToAdd := make(map[string]interface{})
	for _, element := range plan.Privileges.Elements() {
		elementAttr, castOk := element.(types.String)
		if !castOk || elementAttr.IsNull() || elementAttr.IsUnknown() {
			diags.AddAttributeError(
				rolePrivilegesPath.AtSetValue(element),
				"Unable to update Role",
				fmt.Sprintf(
					"Encountered a bad value loading set data: %s",
					element,
				),
			)
			return
		}
		elementValue := elementAttr.ValueString()
		privilegesToAdd[elementValue] = struct{}{}
	}

	// Load all existing privileges, removing from those to be added if found,
	// and added to those to be removed if not.
	privilegesToRemove := make(map[string]interface{})
	for _, element := range plan.Privileges.Elements() {

		elementAttr, castOk := element.(types.String)
		if !castOk || elementAttr.IsNull() || elementAttr.IsUnknown() {
			diags.AddAttributeError(
				rolePrivilegesPath.AtSetValue(element),
				"Unable to update Role",
				fmt.Sprintf(
					"Encountered a bad value loading set data: %s",
					element,
				),
			)
			return
		}
		elementValue := elementAttr.ValueString()

		// If we've found the same privilege in the plan, we can remove it from
		// the privileges we need to add, and don't need to do anything else.
		_, elementFound := privilegesToAdd[elementValue]
		if elementFound {
			delete(privilegesToAdd, elementValue)
			continue
		}

		// If the item isn't found in the plan, it needs to be removed, so gets
		// added to that set instead.
		privilegesToRemove[elementValue] = struct{}{}

	}

	// Add all the privileges that need to be added
	addApiRes, addApiErr := client.AddRolePrivilegesWithResponse(
		ctx,
		plan.Id.ValueString(),
		&v3.AddRolePrivilegesParams{},
		v3.AddRolePrivilegesJSONRequestBody{
			Privileges: Ptr(maps.Keys(privilegesToAdd)),
		},
	)
	if addApiErr != nil {
		diags.AddAttributeError(
			rolePrivilegesPath,
			"Unable to add privileges to role",
			fmt.Sprintf(
				"Api error encountered updating %s: %s",
				plan.Id.ValueString(),
				addApiErr,
			),
		)
		return
	}

	// Handle error responses.
	if addApiRes.StatusCode() != 200 {
		CheckApiErrorV3(diags,
			addApiRes.JSON400,
			addApiRes.JSON401,
			addApiRes.JSON403,
			addApiRes.JSON404,
			addApiRes.JSON500,
			addApiRes.JSON502,
			addApiRes.JSON503,
		)
		if !diags.HasError() {
			errHandler(RequireHttpStatus(200, &addApiRes.ClientResponse))
		}
		return
	}

	// Remove all the privileges that need to be removed
	remApiRes, remApiErr := client.RemoveRolePrivilegesWithResponse(
		ctx,
		plan.Id.ValueString(),
		&v3.RemoveRolePrivilegesParams{},
		v3.RemoveRolePrivilegesJSONRequestBody{
			Privileges: Ptr(maps.Keys(privilegesToRemove)),
		},
	)
	if remApiErr != nil {
		diags.AddAttributeError(
			rolePrivilegesPath,
			"Unable to remove privileges from role",
			fmt.Sprintf(
				"Api error encountered updating %s: %s",
				plan.Id.ValueString(),
				remApiErr,
			),
		)
		return
	}

	// Handle error responses.
	if remApiRes.StatusCode() != 200 {
		CheckApiErrorV3(diags,
			remApiRes.JSON400,
			remApiRes.JSON401,
			remApiRes.JSON403,
			remApiRes.JSON404,
			remApiRes.JSON500,
			remApiRes.JSON502,
			remApiRes.JSON503,
		)
		if !diags.HasError() {
			errHandler(RequireHttpStatus(200, &remApiRes.ClientResponse))
		}
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &plan)...)

}

func (r RoleResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadDelete)

	client := r.GetApiClientV3(diags, MsgResourceBadDelete)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RoleResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	apiRes, apiErr := client.DeleteRoleWithResponse(ctx, data.Id.ValueString(), &v3.DeleteRoleParams{})
	if errHandler(apiErr); diags.HasError() {
		return
	}

	errHandler(RequireHttpStatus(200, &apiRes.ClientResponse))
	if diags.HasError() {
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}
