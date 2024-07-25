package provider

import (
	"context"

	. "github.com/hashicorp/terraform-plugin-framework/resource"
)

type IdmcProviderResource struct {
	*IdmcProviderData
}

func (r IdmcProviderResource) Configure(ctx context.Context, req ConfigureRequest, res *ConfigureResponse) {
	r.IdmcProviderData = GetProviderData(&res.Diagnostics, req.ProviderData)
	res.Diagnostics.AddWarning("Resource Configured", "Just testing if this actually happens.")
}
