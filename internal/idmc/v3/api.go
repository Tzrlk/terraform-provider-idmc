package v3

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-idmc/internal/idmc/common"
	. "terraform-provider-idmc/internal/utils"
)

//go:generate -command oapi-codegen go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
//go:generate oapi-codegen -config ./codegen.yml ./openapi.yml

type IdmcAdminV3Api struct {
	Client         ClientWithResponses
	rolePrivileges *map[string]RolePrivilegeItem
}

func NewIdmcAdminV3Api(baseUrl string, sessionId *string, opts ...common.ClientOption) (*IdmcAdminV3Api, error) {

	// Add a request editor to apply the needed api headers on all requests.
	opts = append(opts, common.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header["Accept"] = []string{"application/json"}
		if sessionId != nil {
			req.Header["INFA-SESSION-ID"] = []string{*sessionId}
		}
		return nil
	}))

	apiClient, clientErr := NewClientWithResponses(baseUrl, opts...)
	if clientErr != nil {
		return nil, clientErr
	}

	return OkPtr(&IdmcAdminV3Api{
		Client: *apiClient,
	})

}

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
	if apiRes.StatusCode() != 200 {
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
			return nil, fmt.Errorf(MsgGetRolePrivilegesFailed,
				fmt.Errorf("received http %s but expected %d", apiRes.Status(), 200))
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
	privNames := TransformSlice(privIds, func(from string) string {
		return privMap[from].Name
	})

	return privNames, nil

}

const MsgAddRolePrivilegesFailed = "unable to add %d privileges to role %s: %v"

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
	if apiRes.StatusCode() == 204 {
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
			fmt.Errorf("received http %s but expected %d", apiRes.Status(), 200))
	}

	return fmt.Errorf(MsgAddRolePrivilegesFailed, len(privIds), roleId,
		fmt.Errorf("received api error response: %s", *errBody))

}

const MsgRemoveRolePrivilegesFailed = "unable to remove %d privileges from role %s: %v"

func (i *IdmcAdminV3Api) RemoveRolePrivileges(ctx context.Context, roleId string, privIds []string) error {

	// Exit early if we don't have anything to do.
	if len(privIds) == 0 {
		return nil
	}

	// Translate all the privilege ids into names.
	privNames, err := i.lookupRolePrivilegeNames(ctx, privIds)
	if err != nil {
		return fmt.Errorf(MsgRemoveRolePrivilegesFailed, len(privIds), roleId, err)
	}

	// Set up the request inputs
	apiParams := &AddRolePrivilegesParams{}
	apiBody := AddRolePrivilegesJSONRequestBody{
		Privileges: privNames,
	}

	apiRes, apiErr := i.Client.AddRolePrivilegesWithResponse(ctx, roleId, apiParams, apiBody)
	if apiErr != nil {
		return fmt.Errorf(MsgRemoveRolePrivilegesFailed, len(privIds), roleId, apiErr)
	}

	// Return early if everything was successful.
	if apiRes.StatusCode() == 204 {
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
		return fmt.Errorf(MsgRemoveRolePrivilegesFailed, len(privIds), roleId,
			fmt.Errorf("received http %s but expected %d", apiRes.Status(), 200))
	}

	return fmt.Errorf(MsgRemoveRolePrivilegesFailed, len(privIds), roleId,
		fmt.Errorf("received api error response: %s", *errBody))

}

func (a ApiErrorResponseBody) String() string {
	var msg strings.Builder
	msg.WriteString("--- error ---\n")
	msg.WriteString(a.Error.String())
	return msg.String()
}

func (a ApiError) String() string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("Request: %s\n", a.RequestId))
	msg.WriteString(fmt.Sprintf("Code:    %s\n", a.Code))
	msg.WriteString(fmt.Sprintf("Msg:     %s\n", a.Message))
	if a.DebugMessage != nil && *a.DebugMessage != a.Message {
		msg.WriteString(fmt.Sprintf("Debug:   %s\n", *a.DebugMessage))
	}
	if a.Details != nil {
		msg.WriteString("Details:\n")
		for _, detail := range *a.Details {
			msg.WriteString(fmt.Sprintf("  - Code:    %s\n", detail.Code))
			msg.WriteString(fmt.Sprintf("    Message: %s\n", detail.Message))
			if detail.DebugMessage != nil && *detail.DebugMessage != detail.Message {
				msg.WriteString(fmt.Sprintf("    Debug:   %s\n", *detail.DebugMessage))
			}
		}
	}
	return msg.String()
}
