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
			"org_id": schema.StringAttribute{
				Description: "Organization ID.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Runtime environment name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the runtime environment.",
				Computed:    true,
			},
			"shared": schema.BoolAttribute{
				Description: "Indicates whether the Secure Agent group is shared.",
				Optional:    true,
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
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadCreate)

	client := r.GetApiClientV2(diags, MsgResourceBadCreate)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.Append(req.Plan.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	reqBody := v2.CreateRuntimeEnvironmentJSONRequestBody{
		Type:     utils.Ptr(v2.RuntimeEnvironmentDataMinimalTypeRuntimeEnvironment),
		Name:     data.Name.ValueString(),
		IsShared: data.Shared.ValueBoolPointer(),
	}

	apiRes, apiErr := client.CreateRuntimeEnvironmentWithResponse(ctx, reqBody)
	if errHandler(apiErr); diags.HasError() {
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
			errHandler(RequireHttpStatus(200, apiRes))
		}
		return
	}

	errHandler(r.UpdateState(diags, &data, apiRes.JSON200))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

// </editor-fold>

// Create <editor-fold desc="Read" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadRead)

	client := r.GetApiClientV2(diags, MsgResourceBadRead)
	if diags.HasError() {
		return
	}

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	if data.Id.IsNull() {
		diags.AddAttributeError(
			path.Root("id"),
			MsgResourceBadRead,
			"Resource id is missing.",
		)
		return
	}

	// Perform the API request.
	apiRes, apiErr := client.GetRuntimeEnvironmentWithResponse(ctx, data.Id.ValueString())
	if errHandler(apiErr); apiErr != nil {
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
			errHandler(RequireHttpStatus(200, apiRes))
		}
		return
	}

	errHandler(r.UpdateState(diags, &data, apiRes.JSON200))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

// </editor-fold>

// Update <editor-fold desc="Update" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadUpdate)

	client := r.GetApiClientV2(diags, MsgResourceBadUpdate)
	if diags.HasError() {
		return
	}

	// Load config from state for comparison.
	var state RuntimeEnvironmentResourceModel
	diags.Append(req.State.Get(ctx, &state)...)

	// Load configuration from plan.
	var plan RuntimeEnvironmentResourceModel
	diags.Append(req.Plan.Get(ctx, &plan)...)

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
	if errHandler(apiErr); diags.HasError() {
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
			errHandler(RequireHttpStatus(200, apiRes))
		}
		return
	}

	errHandler(r.UpdateState(diags, &plan, apiRes.JSON200))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &plan)...)

}

// </editor-fold>

// Delete <editor-fold desc="Delete" defaultstate="collapsed">
func (r RuntimeEnvironmentResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, MsgResourceBadDelete)

	client := r.GetApiClientV2(diags, MsgResourceBadDelete)
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
	if errHandler(apiErr); diags.HasError() {
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
			errHandler(RequireHttpStatus(200, apiRes))
		}
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

// </editor-fold>

func (r RuntimeEnvironmentResource) UpdateState(
	diags *diag.Diagnostics,
	state *RuntimeEnvironmentResourceModel,
	data *v2.RuntimeEnvironment,
) error {
	if data == nil {
		return fmt.Errorf("no runtime environment response data provided")
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

	if data.Agents == nil {
		state.Agents = types.SetUnknown(types.StringType)
	} else {
		agents := make([]attr.Value, len(*data.Agents))
		for index, agent := range *data.Agents {
			agents[index] = types.StringPointerValue(agent.Id)
		}
		state.Agents = UnwrapSetValue(diags, path.Root("agents"), types.StringType, agents)
	}

	return nil
}
