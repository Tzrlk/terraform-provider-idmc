package idmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ApiV3LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApiV3LoginResponse struct {
	Products []ApiV3LoginResponseProduct `json:"products"`
	UserInfo ApiV3LoginResponseUserInfo  `json:"userInfo"`
}

type ApiV3LoginResponseProduct struct {
	Name       string `json:"name"`
	BaseApiUrl string `json:"baseApiUrl"`
}

type ApiV3LoginResponseUserInfo struct {
	SessionId   string            `json:"sessionId"`
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	ParentOrgId string            `json:"parentOrgId"`
	OrgId       string            `json:"orgId"`
	OrgName     string            `json:"orgName"`
	Groups      map[string]string `json:"groups"`
	Status      string            `json:"status"`
}

func (api *ApiV3) DoLogin(diag *diag.Diagnostics, authUser string, authPass string) *Api {

	// Set up the login function config
	authFunction := ApiItem[ApiV3LoginRequest, ApiV3LoginResponse]{
		Api:    api.Root,
		Method: "POST",
		Path:   "core/v3/login",
	}

	// Attempt to perform the login request
	authResponse := authFunction.Call(diag, nil, nil, &ApiV3LoginRequest{
		Username: authUser,
		Password: authPass,
	})
	if diag.HasError() {
		return nil
	}

	// A new explicit api struct is created to allow for a single base api to
	// be used for multiple user logins if needed.
	authedApi := Api{
		Client:    api.Root.Client,
		BaseUrl:   "",
		SessionId: authResponse.UserInfo.SessionId,
	}

	// We're looking for the integration cloud baseurl. If that isn't found, we
	// need to present an error instead of continuing.
	for idx := range authResponse.Products {
		if authResponse.Products[idx].Name == "Integration Cloud" {
			authedApi.BaseUrl = authResponse.Products[idx].BaseApiUrl
			break
		}
	}
	if authedApi.BaseUrl == "" {
		diag.AddError(
			"Unable to determine Base API Url",
			fmt.Sprintf("Integration Cloud not found in login response products: %v", authResponse.Products),
		)
		return nil
	}

	return &authedApi

}
