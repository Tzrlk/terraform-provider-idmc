package idmc

import "terraform-provider-idmc/internal/idmc/admin"

type IdmcApi struct {
	Admin *admin.IdmcAdminApi
}
