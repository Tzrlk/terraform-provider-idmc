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

	// Perform the login operation with the provided credentials.
	baseApiUrl, sessionId, loginError := doLogin(ctx, authHost, authUser, authPass)
	if loginError != nil {
		return nil, loginError
	}

	// Construct the actual client implementations.
	idmcApi := &IdmcApi{}

	idmcAdminApi, idmcAdminApiErr := admin.NewIdmcAdminApi(*baseApiUrl, sessionId)
	if idmcAdminApiErr != nil {
		return nil, idmcAdminApiErr
	}

	idmcApi.Admin = idmcAdminApi
	return idmcApi, nil

}

func doLogin(
	ctx context.Context,
	authHost string,
	authUser string,
	authPass string,
	opts ...v3.ClientOption,
) (*string, *string, error) {

	// First set up a client configured for api login.
	loginServerUrl := fmt.Sprintf("https://%s/public", authHost)
	loginClient, loginClientError := v3.NewClientWithResponses(loginServerUrl, opts...)
	if loginClientError != nil {
		return nil, nil, loginClientError
	}

	// Perform the login operation with the provided credentials.
	loginResponse, loginResponseError := loginClient.LoginWithResponse(ctx, v3.LoginJSONRequestBody{
		Username: authUser,
		Password: authPass,
	})
	if loginResponseError != nil {
		return nil, nil, loginResponseError
	}

	// Extract the key information from the login response
	if loginResponse.StatusCode() != 200 {
		return nil, nil, fmt.Errorf(
			"expected http 200 ok, got %s",
			loginResponse.Status(),
		)
	}
	if loginResponse.JSON200 == nil {
		return nil, nil, fmt.Errorf(
			"expected response to be parsed as json, found nil",
		)
	}
	loginResponseJson := *loginResponse.JSON200
	userInfo := *loginResponseJson.UserInfo

	sessionId := userInfo.SessionId
	for _, product := range *loginResponse.JSON200.Products {
		if *product.Name == "Integration Cloud" {
			return product.BaseApiUrl, sessionId, nil
		}
	}

	// TODO: This should probably just return an error, or fall-back to other products.
	return &loginServerUrl, sessionId, nil

}
