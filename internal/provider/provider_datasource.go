package provider

import (
	"context"
	. "github.com/hashicorp/terraform-plugin-framework/datasource"
)

type IdmcProviderDataSource struct {
	*IdmcProviderData
}

func (d IdmcProviderDataSource) Configure(ctx context.Context, req ConfigureRequest, res *ConfigureResponse) {
	d.IdmcProviderData = GetProviderData(&res.Diagnostics, req.ProviderData)
	res.Diagnostics.AddWarning("Datasource Configured", "Just testing if this actually happens.")
}
