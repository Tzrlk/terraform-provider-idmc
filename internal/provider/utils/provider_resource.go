package utils

import (
	"context"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	MsgResourceBadConfig = "Unable to configure resource"
	MsgResourceBadUpdate = "Unable to read resource"
	MsgResourceBadDelete = "Unable to delete resource"
	MsgResourceBadRead   = "Unable to read resource"
	MsgResourceBadCreate = "Unable to create resource"
)

type IdmcProviderResource struct {
	*IdmcProviderData
}

func (r *IdmcProviderResource) Configure(ctx context.Context, req ConfigureRequest, res *ConfigureResponse) {
	diags := NewDiagsHandler(&res.Diagnostics, MsgResourceBadConfig)
	r.IdmcProviderData = GetProviderData(diags, req.ProviderData)
	if r.IdmcProviderData == nil && req.ProviderData != nil {
		diags.AddError("GetProviderData returned nil, but the original value isn't.")
	}
}
