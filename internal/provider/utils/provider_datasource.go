package utils

import (
	"context"
	. "github.com/hashicorp/terraform-plugin-framework/datasource"
)

const (
	MsgDataSourceBadConfig = "Unable to configure data source"
	MsgDataSourceBadRead   = "Unable to read data source"
)

type IdmcProviderDataSource struct {
	*IdmcProviderData
}

func (d *IdmcProviderDataSource) Configure(ctx context.Context, req ConfigureRequest, res *ConfigureResponse) {
	if d.IdmcProviderData != nil && d.IdmcProviderData.Api != nil {
		return // just leave it.
	}
	diags := NewDiagsHandler(&res.Diagnostics, MsgDataSourceBadConfig)
	d.IdmcProviderData = GetProviderData(diags, req.ProviderData)
	if d.IdmcProviderData == nil && req.ProviderData != nil {
		diags.AddError("GetProviderData returned nil, but the original value isn't.")
	}
}
