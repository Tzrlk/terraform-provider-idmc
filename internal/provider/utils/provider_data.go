package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-idmc/internal/idmc"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"
)

type IdmcProviderData struct {
	Api *idmc.IdmcApi
}

func (r *IdmcProviderData) GetApi(diags *diag.Diagnostics, errTitle string) *idmc.IdmcApi {
	if r == nil {
		diags.AddError(errTitle, "The provider (and therefore IDMC api client) has not been configured yet.")
		return nil
	}
	if r.Api == nil {
		diags.AddError(errTitle, "The provider has not properly initialised the api client.")
	}
	return r.Api
}

func (r *IdmcProviderData) GetApiClientV2(diags *diag.Diagnostics, errTitle string) *v2.ClientWithResponses {
	api := r.GetApi(diags, errTitle)
	if api == nil {
		return nil
	}
	if api.V2 == nil {
		diags.AddError(errTitle, "The V2 api client wrapper has not been initialised.")
		return nil
	}
	if api.V2.Client == nil {
		diags.AddError(errTitle, "The V2 api client has not been initialised.")
		return nil
	}
	return api.V2.Client
}

func (r *IdmcProviderData) GetApiClientV3(diags *diag.Diagnostics, errTitle string) *v3.ClientWithResponses {
	api := r.GetApi(diags, errTitle)
	if api == nil {
		return nil
	}
	if api.V3 == nil {
		diags.AddError(errTitle, "The V3 api client wrapper has not been initialised.")
		return nil
	}
	if api.V3.Client == nil {
		diags.AddError(errTitle, "The V3 api client has not been initialised.")
		return nil
	}
	return api.V3.Client
}

func GetProviderData(diags *diag.Diagnostics, data any, errTitle string) *IdmcProviderData {

	// Provider hasn't been configured yet.
	if data == nil {
		return nil
	}

	// Attempt to cast the data.
	if data, ok := data.(*IdmcProviderData); ok {
		return data
	}

	// Really, this should never happen.
	diags.AddError(errTitle, fmt.Sprintf(
		"Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.",
		data,
	))
	return nil

}
