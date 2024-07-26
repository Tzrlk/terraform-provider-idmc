# by id
data "idmc_role" "by_id" {
  id = var.role_id
}

# by name
data "idmc_role" "by_name" {
  name = var.role_name
}
