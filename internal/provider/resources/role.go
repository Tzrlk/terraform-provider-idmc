package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	. "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
	v3 "terraform-provider-idmc/internal/idmc/admin/v3"
	. "terraform-provider-idmc/internal/utils"
)

var _ Resource = &RoleResource{}

func NewRoleResource() Resource {
	return &RoleResource{}
}

type RoleResource struct {
	Client *v3.ClientWithResponses
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

func (r RoleResource) Metadata(ctx context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r RoleResource) Schema(ctx context.Context, req SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Service generated identifier for the role.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"privileges": schema.SetAttribute{
				MarkdownDescription: "",
				Required:            true,
				ElementType:         types.StringType,
			},
			"org_id": schema.StringAttribute{
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
			"system_role": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"status": schema.StringAttribute{
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
			"created_time": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"updated_time": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
		},
	}
}

func (r RoleResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
	diags := &resp.Diagnostics

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

	apiResp, apiErr := r.Client.CreateRoleWithResponse(ctx, v3.CreateRoleJSONRequestBody{
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Privileges:  &rolePrivilegeValues,
	})
	if apiErr != nil {
		diags.AddError(
			"Unable to create resource",
			fmt.Sprintf("Encountered an error communicating with the api: %s", apiErr),
		)
		return
	}

	// Update the configured state so instabilities can be detected.
	data.Id = types.StringPointerValue(apiResp.JSON200.Id)
	data.Name = types.StringPointerValue(apiResp.JSON200.RoleName)
	data.Description = types.StringPointerValue(apiResp.JSON200.Description)
	data.Privileges = UnwrapDiag(diags, path.Root("privileges"), func() (types.Set, diag.Diagnostics) {
		apiRolePrivilegeItems := *apiResp.JSON200.Privileges
		apiRolePrivilegeAttrs := make([]attr.Value, len(apiRolePrivilegeItems))
		for index, apiRolePrivilegeItem := range apiRolePrivilegeItems {
			apiRolePrivilegeAttrs[index] = types.StringPointerValue(apiRolePrivilegeItem.Id)
		}
		return types.SetValue(types.StringType, apiRolePrivilegeAttrs)
	})

	// Update derived values
	data.OrgId = types.StringPointerValue(apiResp.JSON200.OrgId)
	data.DisplayName = types.StringPointerValue(apiResp.JSON200.DisplayName)
	data.DisplayDescription = types.StringPointerValue(apiResp.JSON200.DisplayDescription)
	data.SystemRole = types.BoolPointerValue(apiResp.JSON200.SystemRole)
	data.Status = types.StringPointerValue(apiResp.JSON200.Status)
	data.CreatedBy = types.StringPointerValue(apiResp.JSON200.CreatedBy)
	data.UpdatedBy = types.StringPointerValue(apiResp.JSON200.UpdatedBy)
	data.CreatedTime = types.StringPointerValue(apiResp.JSON200.CreateTime)
	data.UpdatedTime = types.StringPointerValue(apiResp.JSON200.UpdateTime)

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

func (r RoleResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics

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
	apiResp, apiErr := r.Client.GetRolesWithResponse(ctx, params)
	if apiErr != nil {
		diags.AddError(
			"Http Request Failure",
			fmt.Sprintf("IDMC Api request failure: %s", apiErr),
		)
		return
	}

	apiItems := *apiResp.JSON200
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
	data.Status = types.StringPointerValue(apiItems[0].Status)
	data.CreatedBy = types.StringPointerValue(apiItems[0].CreatedBy)
	data.UpdatedBy = types.StringPointerValue(apiItems[0].UpdatedBy)
	data.CreatedTime = types.StringPointerValue(apiItems[0].CreateTime)
	data.UpdatedTime = types.StringPointerValue(apiItems[0].UpdateTime)

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

func (r RoleResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := &resp.Diagnostics

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
	_, addErr := r.Client.AddRolePrivileges(
		ctx,
		plan.Id.ValueString(),
		&v3.AddRolePrivilegesParams{},
		v3.AddRolePrivilegesJSONRequestBody{
			Privileges: Ptr(maps.Keys(privilegesToAdd)),
		},
	)
	if addErr != nil {
		diags.AddAttributeError(
			rolePrivilegesPath,
			"Unable to add privileges to role",
			fmt.Sprintf(
				"Api error encountered updating %s: %s",
				plan.Id.ValueString(),
				addErr,
			),
		)
		return
	}

	// Remove all the privileges that need to be removed
	_, removeErr := r.Client.RemoveRolePrivileges(
		ctx,
		plan.Id.ValueString(),
		&v3.RemoveRolePrivilegesParams{},
		v3.RemoveRolePrivilegesJSONRequestBody{
			Privileges: Ptr(maps.Keys(privilegesToRemove)),
		},
	)
	if removeErr != nil {
		diags.AddAttributeError(
			rolePrivilegesPath,
			"Unable to remove privileges from role",
			fmt.Sprintf(
				"Api error encountered updating %s: %s",
				plan.Id.ValueString(),
				addErr,
			),
		)
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &plan)...)

}

func (r RoleResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := &resp.Diagnostics

	// Load configuration from plan.
	var data RoleResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	_, apiErr := r.Client.DeleteRoleWithResponse(ctx, data.Id.ValueString(), &v3.DeleteRoleParams{})
	if apiErr != nil {
		diags.AddError(
			"Http Request Failure",
			fmt.Sprintf("IDMC Api request failure: %s", apiErr),
		)
		return
	}

	// Save creation result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}
