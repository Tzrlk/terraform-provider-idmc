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
	d.IdmcProviderData = GetProviderData(&res.Diagnostics, req.ProviderData, MsgDataSourceBadConfig)
	if d.IdmcProviderData == nil && req.ProviderData != nil {
		res.Diagnostics.AddError(MsgDataSourceBadConfig,
			"GetProviderData returned nil, but the original value isn't.")
	}
}
