package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/utils"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ ResourceWithConfigure = &RuntimeEnvironmentResource{}

type RuntimeEnvironmentResource struct {
	*IdmcProviderResource
}

func NewRuntimeEnvironmentResource() Resource {
	return &RuntimeEnvironmentResource{
		&IdmcProviderResource{},
	}
}

type RuntimeEnvironmentResourceModel struct {
	Id          types.String `tfsdk:"id"`
	OrgId       types.String `tfsdk:"org_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedTime types.String `tfsdk:"created_time"`
	UpdatedTime types.String `tfsdk:"updated_time"`
	CreatedBy   types.String `tfsdk:"created_by"`
	UpdatedBy   types.String `tfsdk:"updated_by"`
	Shared      types.Bool   `tfsdk:"shared"`
	FederatedId types.String `tfsdk:"federated_id"`
	Agents      types.Set    `tfsdk:"agents"`
}

// TODO: Implement serverless config.

// Metadata <editor-fold desc="Metadata" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Metadata(ctx context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runtime_environment"
}

// </editor-fold>

// Schema <editor-fold desc="Schema" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Schema(ctx context.Context, req SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-2-resources/runtime_environments.html",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Runtime environment ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Runtime environment name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shared": schema.BoolAttribute{
				Description: "Indicates whether the Secure Agent group is shared.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the runtime environment.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "Organization ID.",
				Computed:    true,
			},
			"federated_id": schema.StringAttribute{
				Description: "Global unique identifier.",
				Computed:    true,
			},
			"agents": schema.SetAttribute{
				Description: "The agents allocated to this runtime environment.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_by": schema.StringAttribute{
				Description: "User who created the runtime environment.",
				Computed:    true,
			},
			"updated_by": schema.StringAttribute{
				Description: "User who last updated the runtime environment.",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Date and time the runtime environment was created.",
				Computed:    true,
			},
			"updated_time": schema.StringAttribute{
				Description: "Date and time that the runtime environment was last updated.",
				Computed:    true,
			},
		},
	}
}

// </editor-fold>

// Create <editor-fold desc="Create" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
	diags := NewDiagsHandler(&resp.Diagnostics, MsgResourceBadCreate)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV2(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.HandleDiags(req.Plan.Get(ctx, &data))
	if diags.HasError() {
		return
	}

	reqBody := v2.CreateRuntimeEnvironmentJSONRequestBody{
		Type:     utils.Ptr(v2.RuntimeEnvironmentDataMinimalTypeRuntimeEnvironment),
		Name:     data.Name.ValueString(),
		IsShared: data.Shared.ValueBoolPointer(),
	}

	apiRes, apiErr := client.CreateRuntimeEnvironmentWithResponse(ctx, reqBody)
	if diags.HandleError(apiErr) {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode() != 200 {
		CheckApiErrorV2(diags,
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

	if r.updateRuntimeEnvironmentState(diags, &data, apiRes.JSON200) {
		return
	}

	// Save result back to state.
	diags.HandleDiags(resp.State.Set(ctx, &data))

}

// </editor-fold>

// Create <editor-fold desc="Read" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := NewDiagsHandler(&resp.Diagnostics, MsgResourceBadRead)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV2(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	if diags.HandleDiags(req.State.Get(ctx, &data)) {
		return
	}

	if data.Id.IsNull() {
		diags.WithPath(path.Root("id")).HandleErrMsg(
			"Resource id is missing.")
		return
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRuntimeEnvironmentWithResponse(ctx, data.Id.ValueString())
	if diags.HandleError(apiErr) {
		return
	}

	// Remove the resource if not found.
	if apiRes.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Handle remaining error responses.
	if apiRes.StatusCode() != 200 {
		CheckApiErrorV2(diags,
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

	if r.updateRuntimeEnvironmentState(diags, &data, apiRes.JSON200) {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

// </editor-fold>

// Update <editor-fold desc="Update" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := NewDiagsHandler(&resp.Diagnostics, MsgResourceBadUpdate)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV2(diags)
	if diags.HasError() {
		return
	}

	// Load config from state for comparison.
	var state RuntimeEnvironmentResourceModel
	diags.HandleDiags(req.State.Get(ctx, &state))

	// Load configuration from plan.
	var plan RuntimeEnvironmentResourceModel
	diags.HandleDiags(req.Plan.Get(ctx, &plan))

	// Only check for errors here so we can see if there are any issues with
	// either data structure before breaking.
	if diags.HasError() {
		return
	}

	// Convert the stored agent ids to keep things consistent.
	agents := make([]v2.RuntimeEnvironmentAgent, len(plan.Agents.Elements()))
	for index, element := range plan.Agents.Elements() {
		if stringVal, ok := element.(types.String); ok {
			agents[index] = v2.RuntimeEnvironmentAgent{
				Id:    stringVal.ValueStringPointer(),
				OrgId: plan.OrgId.ValueStringPointer(),
			}
		}
	}

	reqBody := v2.UpdateRuntimeEnvironmentJSONRequestBody{
		Name:     plan.Name.ValueString(),
		IsShared: plan.Shared.ValueBoolPointer(),
		Agents:   &agents,
	}

	apiRes, apiErr := client.UpdateRuntimeEnvironmentWithResponse(ctx, plan.Id.ValueString(), reqBody)
	if diags.HandleError(apiErr) {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode() != 200 {
		CheckApiErrorV2(diags,
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

	if r.updateRuntimeEnvironmentState(diags, &plan, apiRes.JSON200) {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &plan)...)

}

// </editor-fold>

// Delete <editor-fold desc="Delete" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := NewDiagsHandler(&resp.Diagnostics, MsgResourceBadDelete)
	defer func() { diags.HandlePanic(recover()) }()

	client := r.GetApiClientV2(diags)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	apiRes, apiErr := client.DeleteRuntimeEnvironmentWithResponse(ctx, data.Id.ValueString())
	if diags.HandleError(apiErr) {
		return
	}

	// Handle error responses.
	if apiRes.StatusCode() != 200 {
		CheckApiErrorV2(diags,
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

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

// </editor-fold>

func (r RuntimeEnvironmentResource) updateRuntimeEnvironmentState(
	diags DiagsHandler,
	state *RuntimeEnvironmentResourceModel,
	data *v2.RuntimeEnvironment,
) bool {
	if data == nil {
		diags.HandleErrMsg("no runtime environment response data provided")
		return true
	}

	// Update the configured state so instabilities can be detected.
	state.Id = types.StringPointerValue(data.Id)
	state.Name = types.StringValue(data.Name)
	state.Description = types.StringPointerValue(data.Description)
	state.Shared = types.BoolPointerValue(data.IsShared)

	// Update derived values
	state.OrgId = types.StringPointerValue(data.OrgId)
	state.CreatedBy = types.StringPointerValue(data.CreatedBy)
	state.UpdatedBy = types.StringPointerValue(data.UpdatedBy)
	state.CreatedTime = types.StringPointerValue(data.CreateTime)
	state.UpdatedTime = types.StringPointerValue(data.UpdateTime)
	state.FederatedId = types.StringPointerValue(data.FederatedId)

	agentsDiags := diags.AtName("agents")
	if data.Agents == nil {
		diags.WithTitle("Issue handling API response").HandleWarnMsg(
			"Runtime Environment is expected to have at least an empty list of agents.")
		state.Agents = types.SetNull(types.StringType)
		return diags.HasError()
	}

	agentsAttrs := make([]attr.Value, len(*data.Agents))
	for index, agent := range *data.Agents {
		agentsAttrs[index] = types.StringPointerValue(agent.Id)
	}

	agentsAttr := agentsDiags.SetValue(types.StringType, agentsAttrs)
	if diags.HasError() {
		state.Agents = types.SetUnknown(types.StringType)
		return true
	}

	state.Agents = agentsAttr
	return diags.HasError()

}
