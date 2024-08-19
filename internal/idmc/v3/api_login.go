package v3

import (
	"context"
	"fmt"
	"net/http"

	. "terraform-provider-idmc/internal/utils"
)

const msgLoginFailed = "failed to log-in to %s as %s: %w"

func (i *IdmcAdminV3Api) Login(ctx context.Context, authUser string, authPass string) error {

	// Set up request inputs.
	reqBody := LoginJSONRequestBody{
		Username: authUser,
		Password: authPass,
	}

	// Perform the login operation with the provided credentials.
	res, err := i.Client.LoginWithResponse(ctx, reqBody)
	if err != nil {
		return fmt.Errorf(msgLoginFailed, i.Client.Server, authUser, err)
	}

	// We only want 200 responses.
	if res.StatusCode != http.StatusOK {

		errBody := Coalesce(
			res.JSON400,
			res.JSON401,
			res.JSON403,
			res.JSON404,
			res.JSON500,
			res.JSON502,
			res.JSON503,
		)
		if errBody != nil {
			return fmt.Errorf(msgLoginFailed, i.Client.Server, authUser,
				fmt.Errorf("received api error response: %s", *errBody))
		}

		return fmt.Errorf(msgLoginFailed, i.Client.Server, authUser,
			fmt.Errorf("expected 200 OK, but received %s", res.Status))

	}

	// Check that we actually received a login response payload.
	if res.JSON200 == nil {
		return fmt.Errorf(msgLoginFailed, i.Client.Server, authUser,
			fmt.Errorf("data from response body not parsed"))
	}

	// Update the client session id to the one just received.
	i.Client.SessionId = res.JSON200.UserInfo.SessionId

	// Attempt to update the client url
	for _, product := range res.JSON200.Products {
		if product.Name == "Integration Cloud" {
			if err = i.Client.SetServer(product.BaseApiUrl); err != nil {
				return fmt.Errorf(msgLoginFailed, i.Client.Server, authUser,
					fmt.Errorf("found new api url, but failed update: %v", err))
			}
			return nil
		}
	}

	// Didn't find integration cloud api url in products.
	return fmt.Errorf(msgLoginFailed, i.Client.Server, authUser,
		fmt.Errorf("no IDMC api url found in login response"))

}
