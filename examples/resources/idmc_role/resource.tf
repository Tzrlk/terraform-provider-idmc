resource "idmc_role" "example" {
  name        = var.role_name
  description = var.role_description
  privileges  = var.role_privileges
}

# Inputs
variable "role_name" {
  type = string
}
variable "role_description" {
  type = string
}
variable "role_privileges" {
  type = list(string)
}

# Outputs
output "example" {
  value = idmc_role.example
}
