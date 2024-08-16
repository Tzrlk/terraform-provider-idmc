package utils

import (
	"terraform-provider-idmc/internal/idmc"
)

type IdmcProviderData struct {
	Api *idmc.IdmcApi
}

func (r *IdmcProviderData) GetApi() idmc.IdmcApi {
	if r == nil {
		panic("the provider has not been configured yet")
	}
	if r.Api == nil {
		panic("the provider has not properly initialised the Api client")
	}
	return *r.Api
}

func GetProviderData(diags DiagsHandler, data any) *IdmcProviderData {

	// Provider hasn't been configured yet.
	if data == nil {
		return nil
	}

	// Attempt to cast the data.
	if data, ok := data.(*IdmcProviderData); ok {
		return data
	}

	// Really, this should never happen.
	diags.AddError(
		"Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.",
		data,
	)
	return nil

}
