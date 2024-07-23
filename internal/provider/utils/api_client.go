package utils

import (
	"fmt"
	"terraform-provider-idmc/internal/idmc/common"
)

func RequireHttpStatus(status int, apiRes common.ApiResponse) error {
	if apiRes.StatusCode() != status {
		return fmt.Errorf("expected http status %d but got %s", status, apiRes.Status())
	}
	return nil
}
