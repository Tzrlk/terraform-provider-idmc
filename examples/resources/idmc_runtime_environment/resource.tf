resource "idmc_runtime_environment" "example" {
  name   = var.env_name
  shared = false
}

# Inputs
variable "env_name" {
  type    = string
}
variable "env_shared" {
  type    = bool
}

# Outputs
output "example" {
  value = idmc_runtime_environment.example
}
