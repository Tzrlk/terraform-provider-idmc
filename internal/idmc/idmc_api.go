package idmc

import (
	"context"
	"fmt"
	"terraform-provider-idmc/internal/idmc/admin"
	v3 "terraform-provider-idmc/internal/idmc/admin/v3"
)

type IdmcApi struct {
	Admin *admin.IdmcAdminApi
}

func NewIdmcApi(
	ctx context.Context,
	authHost string,
	authUser string,
	authPass string,
) (*IdmcApi, error) {

	// First set up a client configured for api login.
	loginServerUrl := fmt.Sprintf("https://%s/public", authHost)
	loginClient, loginClientError := v3.NewClientWithResponses(loginServerUrl)
	if loginClientError != nil {
		return nil, loginClientError
	}

	// Perform the login operation with the provided credentials.
	loginResponse, loginResponseError := loginClient.LoginWithResponse(ctx, v3.LoginJSONRequestBody{
		Username: authUser,
		Password: authPass,
	})
	if loginResponseError != nil {
		return nil, loginResponseError
	}

	// Extract the key information from the login response
	sessionId := loginResponse.JSON200.UserInfo.SessionId
	var apiUrl string
	for _, product := range *loginResponse.JSON200.Products {
		if *product.Name == "" {
			apiUrl = *product.BaseApiUrl
		}
	}

	// Construct the actual client implementations.
	idmcApi := &IdmcApi{}

	idmcAdminApi, idmcAdminApiErr := admin.NewIdmcAdminApi(apiUrl, sessionId)
	if idmcAdminApiErr != nil {
		return nil, idmcAdminApiErr
	}

	idmcApi.Admin = idmcAdminApi
	return idmcApi, nil

}
