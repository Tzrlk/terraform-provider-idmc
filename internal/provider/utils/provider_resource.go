package utils

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/path"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	MsgResourceBadConfig = "Unable to configure resource"
	MsgResourceBadUpdate = "Unable to read resource"
	MsgResourceBadDelete = "Unable to delete resource"
	MsgResourceBadRead   = "Unable to read resource"
	MsgResourceBadCreate = "Unable to create resource"
)

// IdmcProviderResource wraps IdmcProviderData and partially implements the Resource interface to support common
// functionality between all resources.
type IdmcProviderResource struct {
	*IdmcProviderData
	Name string
}

func (r *IdmcProviderResource) Metadata(_ context.Context, req MetadataRequest, resp *MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.Name
}

func (r *IdmcProviderResource) Configure(ctx context.Context, req ConfigureRequest, res *ConfigureResponse) {
	diags := NewDiagsHandler(ctx, &res.Diagnostics, MsgResourceBadConfig)
	r.IdmcProviderData = GetProviderData(diags, req.ProviderData)
	if r.IdmcProviderData == nil && req.ProviderData != nil {
		diags.AddError("GetProviderData returned nil, but the original value isn't.")
	}
}

// ImportState implements ResourceWithImportState for all resources unless overridden.
func (r *IdmcProviderResource) ImportState(ctx context.Context, req ImportStateRequest, res *ImportStateResponse) {
	ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

// Schema configures a bunch of common attributes used across all resources.
func (r *IdmcProviderResource) Schema(_ context.Context, _ SchemaRequest, resp *SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Description: "Service generated resource identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "ID of the organization the resource belongs to.",
				Computed:    true,
			},

			// These will almost certainly need to be overridden, but here
			// because they're pretty standard to exist.
			"name": schema.StringAttribute{
				Description: "Name of the resource.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the resource.",
				Optional:    true,
			},

			"created_by": schema.StringAttribute{
				Description: "User who created the resource.",
				Computed:    true,
			},
			"updated_by": schema.StringAttribute{
				Description: "User who last updated the resource.",
				Computed:    true,
			},
			"created_time": schema.StringAttribute{
				Description: "Date and time the resource was created.",
				CustomType:  timetypes.RFC3339Type{},
				Computed:    true,
			},
			"updated_time": schema.StringAttribute{
				Description: "Date and time the resource was last updated.",
				CustomType:  timetypes.RFC3339Type{},
				Computed:    true,
			},
		},
	}
}
