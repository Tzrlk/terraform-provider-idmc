resource "idmc_runtime_environment" "example" {
  name   = var.name
  shared = var.shared
}

# Inputs
variable "name" {
  type = string
}
variable "shared" {
  type = bool
}

# Outputs
output "example" {
  value = idmc_runtime_environment.example
}
