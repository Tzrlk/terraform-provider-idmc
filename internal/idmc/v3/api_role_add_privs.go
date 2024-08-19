package v3

import (
	"context"
	"fmt"

	. "terraform-provider-idmc/internal/utils"
)

const MsgAddRolePrivilegesFailed = "unable to add %d privileges to role %s: %w"

func (i *IdmcAdminV3Api) AddRolePrivileges(ctx context.Context, roleId string, privIds []string) error {

	// Exit early if we don't have anything to do.
	if len(privIds) == 0 {
		return nil
	}

	// Translate all the privilege ids into names.
	privNames, err := i.lookupRolePrivilegeNames(ctx, privIds)
	if err != nil {
		return fmt.Errorf(MsgAddRolePrivilegesFailed, len(privIds), roleId, err)
	}

	// Set up the request inputs
	apiParams := &AddRolePrivilegesParams{}
	apiBody := AddRolePrivilegesJSONRequestBody{
		Privileges: privNames,
	}

	apiRes, apiErr := i.Client.AddRolePrivilegesWithResponse(ctx, roleId, apiParams, apiBody)
	if apiErr != nil {
		return fmt.Errorf(MsgAddRolePrivilegesFailed, len(privIds), roleId, apiErr)
	}

	// Return early if everything was successful.
	if apiRes.StatusCode == 204 {
		return nil
	}

	errBody := Coalesce(
		apiRes.JSON400,
		apiRes.JSON401,
		apiRes.JSON403,
		apiRes.JSON404,
		apiRes.JSON500,
		apiRes.JSON502,
		apiRes.JSON503,
	)

	if errBody == nil {
		return fmt.Errorf(MsgAddRolePrivilegesFailed, len(privIds), roleId,
			fmt.Errorf("received http %s but expected %d", apiRes.Status, 200))
	}

	return fmt.Errorf(MsgAddRolePrivilegesFailed, len(privIds), roleId,
		fmt.Errorf("received api error response: %s", *errBody))

}
