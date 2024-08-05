package utils

import (
	"fmt"
	"terraform-provider-idmc/internal/idmc"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"
)

type IdmcProviderData struct {
	Api *idmc.IdmcApi
}

func (r *IdmcProviderData) GetApi(diags DiagsHandler) *idmc.IdmcApi {
	if r == nil {
		diags.HandleErrMsg("The provider (and therefore IDMC api client) has not been configured yet.")
		return nil
	}
	if r.Api == nil {
		diags.HandleErrMsg("The provider has not properly initialised the api client.")
	}
	return r.Api
}

func (r *IdmcProviderData) GetApiClientV2(diags DiagsHandler) *v2.ClientWithResponses {
	api := r.GetApi(diags)
	if api == nil {
		return nil
	}
	if api.V2 == nil {
		diags.HandleErrMsg("The V2 api client wrapper has not been initialised.")
		return nil
	}
	if api.V2.Client == nil {
		diags.HandleErrMsg("The V2 api client has not been initialised.")
		return nil
	}
	return api.V2.Client
}

func (r *IdmcProviderData) GetApiClientV3(diags DiagsHandler) *v3.ClientWithResponses {
	api := r.GetApi(diags)
	if api == nil {
		return nil
	}
	if api.V3 == nil {
		diags.HandleErrMsg("The V3 api client wrapper has not been initialised.")
		return nil
	}
	if api.V3.Client == nil {
		diags.HandleErrMsg("The V3 api client has not been initialised.")
		return nil
	}
	return api.V3.Client
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
	diags.HandleErrMsg(fmt.Sprintf(
		"Expected *IdmcProviderData, got: %T. Please report this issue to the provider developers.",
		data,
	))
	return nil

}
