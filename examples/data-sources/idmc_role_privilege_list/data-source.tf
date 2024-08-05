data "idmc_role_privilege_list" "example" {
  status = var.status
}

# Inputs
variable "status" {
  type     = string
  nullable = true
}

# Outputs
output "privileges" {
  value = data.idmc_role_privilege_list.example.privileges
}
