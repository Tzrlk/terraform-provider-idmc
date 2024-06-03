package idmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ApiItemLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApiItemLoginResponse struct {
	Products []ApiItemLoginResponseProduct `json:"products"`
	UserInfo ApiItemLoginResponseUserInfo  `json:"userInfo"`
}

type ApiItemLoginResponseProduct struct {
	Name       string `json:"name"`
	BaseApiUrl string `json:"baseApiUrl"`
}

type ApiItemLoginResponseUserInfo struct {
	SessionId   string            `json:"sessionId"`
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	ParentOrgId string            `json:"parentOrgId"`
	OrgId       string            `json:"orgId"`
	OrgName     string            `json:"orgName"`
	Groups      map[string]string `json:"groups"`
	Status      string            `json:"status"`
}

func (api *Api) DoLogin(diag *diag.Diagnostics, authUser string, authPass string) *Api {

	// Set up the login function config
	authFunction := ApiItem[ApiItemLoginRequest, ApiItemLoginResponse]{
		Api:    api,
		Method: "POST",
		Path:   "core/v3/login",
	}

	// Attempt to perform the login request
	authResponse := authFunction.Call(diag, nil, nil, &ApiItemLoginRequest{
		Username: authUser,
		Password: authPass,
	})
	if diag.HasError() {
		return nil
	}

	// A new explicit api struct is created to allow for a single base api to
	// be used for multiple user logins if needed.
	authedApi := Api{
		Client:    api.Client,
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
