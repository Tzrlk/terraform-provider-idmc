package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-idmc/internal/idmc/common"
	"terraform-provider-idmc/internal/idmc/v2"
	"terraform-provider-idmc/internal/idmc/v3"
	"terraform-provider-idmc/internal/utils"
)

func RequireHttpStatus(status int, apiRes common.ApiResponse) error {
	if apiRes.StatusCode() != status {
		return fmt.Errorf("expected http status %d but got %s", status, apiRes.Status())
	}
	return nil
}

func CheckApiErrorV2(diags *diag.Diagnostics, errs ...*v2.ApiErrorBody) {
	err := utils.Coalesce(errs...)
	if err == nil {
		return
	}
	diags.AddError(
		"IDMC API bad response",
		fmt.Sprintf("Code: %s\nMsg: %s",
			utils.ValOr(err.Code, "-"),
			utils.ValOr(err.Description, "-"),
		),
	)
}

func CheckApiErrorV3(diags *diag.Diagnostics, errs ...*v3.ApiErrorBody) {
	err := utils.Coalesce(errs...)
	if err == nil {
		return
	}
	diags.AddError(
		"IDMC API bad response",
		fmt.Sprintf("Request: %s\nCode: %s\nMsg: %s",
			utils.ValOr(err.Error.RequestId, "-"),
			utils.ValOr(err.Error.Code, "-"),
			utils.ValOr(err.Error.Message, "-"),
		),
	)
}
