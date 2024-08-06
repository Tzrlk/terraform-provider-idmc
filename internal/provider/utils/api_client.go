package utils

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"

	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"

	. "terraform-provider-idmc/internal/utils"
)

const (
	MsgFieldMissing   = "<missing>"
	MsgApiBadResponse = "IDMC API bad response"
)

func RequireHttpStatus(apiRes *common.ClientResponse, statuses ...int) error {
	if !slices.Contains(statuses, apiRes.HTTPResponse.StatusCode) {
		return fmt.Errorf("received http %s but expected %+q", apiRes.HTTPResponse.Status, statuses)
	}
	return nil
}

func CheckApiErrorV2(diags DiagsHandler, apiErrors ...*v2.ApiErrorResponse) {
	diags = diags.WithTitle(MsgApiBadResponse)

	apiError := Coalesce(apiErrors...)
	if apiError == nil {
		return
	}

	var msg strings.Builder
	msg.WriteString("--- error ---")

	// Try handling the error as a v2 error first.
	if v2Error, err := apiError.AsApiErrorResponseBody(); err == nil {
		if v2Error.Type == v2.ApiErrorResponseBodyTypeError {
			msg.WriteString("\nCode: " + v2Error.Code)
			msg.WriteString("\nMsg:  " + v2Error.Description)
			diags.AddError(msg.String())
			return
		}
	}

	// Then try handling it as a v3 formatted error.
	if v3Error, err := apiError.AsExternalRef1ApiErrorResponseBody(); err == nil {
		CheckApiErrorV3(diags, &v3Error)
		return
	}

	// If neither work, just yeet the body in as a string.
	if jsonError, err := apiError.MarshalJSON(); err != nil {
		msg.WriteString("\n")
		msg.Write(jsonError)
		diags.AddError(msg.String())
		return
	}

	msg.WriteString("FAILED TO PARSE")
	diags.AddError(msg.String())
}

func CheckApiErrorV3(diags DiagsHandler, apiErrors ...*v3.ApiErrorResponseBody) {
	diags = diags.WithTitle(MsgApiBadResponse)

	apiError := Coalesce(apiErrors...)
	if apiError == nil {
		return
	}

	var msg strings.Builder
	msg.WriteString("--- error ---")

	msg.WriteString("\nRequest: " + apiError.Error.RequestId)
	msg.WriteString("\nCode:    " + apiError.Error.Code)
	msg.WriteString("\nMsg:     " + apiError.Error.Message)
	if apiError.Error.Details != nil {
		msg.WriteString("\nDetails:")
		for _, detail := range *apiError.Error.Details {
			msg.WriteString("\n  - Code: " + detail.Code)
			msg.WriteString("\n    Msg:  " + detail.Message)
		}
	}

	diags.AddError(msg.String())
}
