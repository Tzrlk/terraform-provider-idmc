package admin

import (
	"terraform-provider-idmc/internal/idmc/admin/v2"
	"terraform-provider-idmc/internal/idmc/admin/v3"
)

type IdmcAdminApi struct {
	V2 *v2.IdmcAdminV2Api
	V3 *v3.IdmcAdminV3Api
}
