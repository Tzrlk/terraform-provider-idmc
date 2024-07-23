package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-idmc/internal/idmc/v2"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
	. "terraform-provider-idmc/internal/provider/utils"
)

var _ Resource = &RuntimeEnvironmentResource{}

func NewRuntimeEnvironmentResource() Resource {
	return &RuntimeEnvironmentResource{}
}

type RuntimeEnvironmentResource struct {
	Client *v2.ClientWithResponses
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
}

// TODO: Implement serverless config.

func (r RuntimeEnvironmentResource) Metadata(ctx context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runtime_environment"
}

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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the runtime environment.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shared": schema.BoolAttribute{
				Description: "Indicates whether the Secure Agent group is shared.",
				Computed:    false,
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

const RuntimeEnvironmentResourceBadCreate = "Unable to create resource"

func (r RuntimeEnvironmentResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, RuntimeEnvironmentResourceBadCreate)

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.Append(req.Plan.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	reqBody := v2.CreateRuntimeEnvironmentJSONRequestBody{
		Name:     data.Name.ValueString(),
		IsShared: data.Shared.ValueBoolPointer(),
	}

	apiRes, apiErr := r.Client.CreateRuntimeEnvironmentWithResponse(ctx, reqBody)
	if errHandler(apiErr); diags.HasError() {
		return
	}

	errHandler(RequireHttpStatus(200, apiRes))
	if diags.HasError() {
		return
	}

	errHandler(r.UpdateState(&data, apiRes.JSON200))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

const RuntimeEnvironmentResourceBadRead = "Unable to read resource"

func (r RuntimeEnvironmentResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, RuntimeEnvironmentResourceBadRead)

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	if data.Id.IsNull() {
		diags.AddAttributeError(
			path.Root("id"),
			RuntimeEnvironmentResourceBadRead,
			"Resource id is missing.",
		)
		return
	}

	// Perform the API request.
	apiRes, apiErr := r.Client.GetRuntimeEnvironmentWithResponse(ctx, data.Id.ValueString())
	if errHandler(apiErr); apiErr != nil {
		return
	}

	errHandler(RequireHttpStatus(200, apiRes))
	if diags.HasError() {
		return
	}

	errHandler(r.UpdateState(&data, apiRes.JSON200))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

const RuntimeEnvironmentResourceBadUpdate = "Unable to read resource"

func (r RuntimeEnvironmentResource) Update(ctx context.Context, req UpdateRequest, resp *UpdateResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, RuntimeEnvironmentResourceBadUpdate)

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

	reqBody := v2.UpdateRuntimeEnvironmentJSONRequestBody{
		Name:     plan.Name.ValueString(), // TODO: will this even work? double-check api.
		IsShared: plan.Shared.ValueBoolPointer(),
		Agents:   nil, // TODO: do this somehow? Skip it?
	}

	apiRes, apiErr := r.Client.UpdateRuntimeEnvironmentWithResponse(ctx, plan.Id.ValueString(), reqBody)
	if errHandler(apiErr); diags.HasError() {
		return
	}

	errHandler(RequireHttpStatus(200, apiRes))
	if diags.HasError() {
		return
	}

	errHandler(r.UpdateState(&plan, apiRes.JSON200))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &plan)...)

}

const RuntimeEnvironmentResourceBadDelete = "Unable to delete resource"

func (r RuntimeEnvironmentResource) Delete(ctx context.Context, req DeleteRequest, resp *DeleteResponse) {
	diags := &resp.Diagnostics
	errHandler := DiagsErrHandler(diags, RuntimeEnvironmentResourceBadDelete)

	// Load configuration from plan.
	var data RuntimeEnvironmentResourceModel
	diags.Append(req.State.Get(ctx, &data)...)
	if diags.HasError() {
		return
	}

	apiRes, apiErr := r.Client.DeleteRuntimeEnvironmentWithResponse(ctx, data.Id.ValueString())
	if errHandler(apiErr); diags.HasError() {
		return
	}

	errHandler(RequireHttpStatus(200, apiRes))
	if diags.HasError() {
		return
	}

	// Save result back to state.
	diags.Append(resp.State.Set(ctx, &data)...)

}

func (r RuntimeEnvironmentResource) UpdateState(
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

	return nil
}
