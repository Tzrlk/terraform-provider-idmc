package utils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"

	. "terraform-provider-idmc/internal/utils"
)

const (
	MsgFieldMissing   = "<missing>"
	MsgApiBadResponse = "IDMC API bad response"
)

func RequireHttpStatus(status int, apiRes *common.ClientResponse) error {
	if apiRes.HTTPResponse.StatusCode != status {
		return fmt.Errorf("expected http status %d but got %s", status, apiRes.HTTPResponse.Status)
	}
	return nil
}

func CheckApiErrorV2(diags *diag.Diagnostics, apiErrors ...*v2.ApiErrorResponse) {
	apiError := Coalesce(apiErrors...)
	if apiError == nil {
		return
	}
	var msg strings.Builder
	msg.WriteString("--- error ---")

	// Try handling the error as a v2 error first.
	v2Error, err := apiError.AsApiErrorResponseBody()
	if err == nil && v2Error.Type == v2.ApiErrorResponseBodyTypeError {
		msg.WriteString("\nCode: " + v2Error.Code)
		msg.WriteString("\nMsg:  " + v2Error.Description)
		diags.AddError(MsgApiBadResponse, msg.String())
		return
	}

	// Then try handling it as a v3 formatted error.
	v3Error, err := apiError.AsV3ApiErrorResponseBody()
	if err == nil {
		msg.WriteString("\nRequest: " + v3Error.Error.RequestId)
		msg.WriteString("\nCode:    " + v3Error.Error.Code)
		msg.WriteString("\nMsg:     " + v3Error.Error.Message)
		if v3Error.Error.Details != nil {
			msg.WriteString("\nDetails:")
			for _, detail := range *v3Error.Error.Details {
				msg.WriteString("\n  - Code: " + detail.Code)
				msg.WriteString("\n    Msg:  " + detail.Message)
			}
		}
		diags.AddError(MsgApiBadResponse, msg.String())
		return
	}

	// If neither work, just yeet the body in as a string.
	jsonError, err := apiError.MarshalJSON()
	if err == nil {
		msg.WriteString("\n")
		msg.Write(jsonError)
		diags.AddError(MsgApiBadResponse, msg.String())
		return
	}

	msg.WriteString("FAILED TO PARSE")
	diags.AddError(MsgApiBadResponse, msg.String())
}

func CheckApiErrorV3(diags *diag.Diagnostics, apiErrors ...*v3.ApiErrorResponseBody) {
	apiError := Coalesce(apiErrors...)
	if apiError == nil {
		return
	}
	var msg strings.Builder
	msg.WriteString("--- error ---")
	msg.WriteString("\nRequest: " + ValOr(apiError.Error.RequestId, MsgFieldMissing))
	msg.WriteString("\nCode:    " + ValOr(apiError.Error.Code, MsgFieldMissing))
	msg.WriteString("\nMsg:     " + ValOr(apiError.Error.Message, MsgFieldMissing))
	if apiError.Error.Details != nil {
		msg.WriteString("\nDetails:")
		for _, detail := range *apiError.Error.Details {
			msg.WriteString("\n  - Code: " + ValOr(detail.Code, MsgFieldMissing))
			msg.WriteString("\n    Msg:  " + ValOr(detail.Message, MsgFieldMissing))
		}
	}
	diags.AddError("IDMC API bad response", msg.String())
}
