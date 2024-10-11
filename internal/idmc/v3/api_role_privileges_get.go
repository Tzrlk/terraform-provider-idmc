package v3

import (
	"context"
	"fmt"
	"github.com/samber/lo"
)

const MsgGetRolePrivilegesFailed = "failed to fetch role privileges: %v"

func (i *IdmcAdminV3Api) GetRolePrivileges(ctx context.Context, status *string) (map[string]RolePrivilegeItem, error) {

	// Used whatever is cached when dealing with the default status query.
	if i.rolePrivileges != nil && status != nil {
		return *i.rolePrivileges, nil
	}

	// Obtain request parameters from config.
	params := &ListPrivilegesParams{}
	if status != nil {
		query := fmt.Sprintf("status==\"%s\"", *status)
		params.Q = &query
	}

	// Fetch the privilege list to map ids to names.
	// TODO: get the _whole_ list (q=status=="All") .
	apiRes, apiErr := i.Client.ListPrivilegesWithResponse(ctx, &ListPrivilegesParams{})
	if apiErr != nil {
		return nil, fmt.Errorf(MsgGetRolePrivilegesFailed,
			apiErr)
	}

	// Handle error responses.
	if apiRes.StatusCode != 200 {
		errBody := lo.CoalesceOrEmpty(
			apiRes.JSON400,
			apiRes.JSON401,
			apiRes.JSON403,
			apiRes.JSON404,
			apiRes.JSON500,
			apiRes.JSON502,
			apiRes.JSON503,
		)

		if errBody == nil {
			return nil, fmt.Errorf(MsgGetRolePrivilegesFailed,
				fmt.Errorf("received http %s but expected %d", apiRes.Status, 200))
		}

		return nil, fmt.Errorf(MsgGetRolePrivilegesFailed,
			fmt.Errorf("received api error response: %s", *errBody))

	}

	if apiRes.JSON200 == nil {
		return nil, fmt.Errorf(MsgGetRolePrivilegesFailed,
			fmt.Errorf("received no valid response payload"))
	}

	rolePrivileges := make(map[string]RolePrivilegeItem, len(*apiRes.JSON200))
	for _, item := range *apiRes.JSON200 {
		rolePrivileges[item.Id] = item
	}

	// Only cache the results of a default query.
	if status == nil {
		i.rolePrivileges = &rolePrivileges
	}

	return rolePrivileges, nil

}

func (i *IdmcAdminV3Api) lookupRolePrivilegeNames(ctx context.Context, privIds []string) ([]string, error) {

	// Exit early if we don't have anything to do.
	if len(privIds) == 0 {
		return make([]string, 0), nil
	}

	// Fetch the role privilege list so we can reference it.
	privMap, err := i.GetRolePrivileges(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Translate all the privilege ids into names.
	privNames := lo.Map(privIds, func(from string, _index int) string {
		return privMap[from].Name
	})

	return privNames, nil

}
