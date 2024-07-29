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

const MsgFieldMissing = "<missing>"

func RequireHttpStatus(status int, apiRes common.ApiResponse) error {
	if apiRes.StatusCode() != status {
		return fmt.Errorf("expected http status %d but got %s", status, apiRes.Status())
	}
	return nil
}

func CheckApiErrorV2(diags *diag.Diagnostics, errs ...*v2.ApiErrorBody) {
	err := Coalesce(errs...)
	if err == nil {
		return
	}
	diags.AddError(
		"IDMC API bad response",
		fmt.Sprintf("Code: %s\nMsg: %s",
			ValOr(err.Code, "-"),
			ValOr(err.Description, "-"),
		),
	)
}

func CheckApiErrorV3(diags *diag.Diagnostics, errs ...*v3.ApiErrorResponseBody) {
	err := Coalesce(errs...)
	if err == nil {
		return
	}
	var msg strings.Builder
	msg.WriteString("--- error ---")
	msg.WriteString("\nRequest: " + ValOr(err.Error.RequestId, MsgFieldMissing))
	msg.WriteString("\nCode:    " + ValOr(err.Error.Code, MsgFieldMissing))
	msg.WriteString("\nMsg:     " + ValOr(err.Error.Message, MsgFieldMissing))
	if err.Error.Details != nil {
		msg.WriteString("\nDetails:")
		for _, detail := range *err.Error.Details {
			msg.WriteString("\n  - Code: " + ValOr(detail.Code, MsgFieldMissing))
			msg.WriteString("\n    Msg:  " + ValOr(detail.Message, MsgFieldMissing))
		}
	}
	diags.AddError("IDMC API bad response", msg.String())
}
