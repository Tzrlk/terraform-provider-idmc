resource "idmc_rbac_role" "example" {
  name        = var.role_name
  description = var.role_description
  privileges = [
    var.role_privilege_id,
  ]
}
