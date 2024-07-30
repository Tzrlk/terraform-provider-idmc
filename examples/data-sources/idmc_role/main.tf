# The correct provider source needs to be selected.
terraform {
  required_providers {
    idmc = {
      source = "tzrlk/idmc"
    }
  }
}

# Needed so we can override it with actual credentials.
provider "idmc" {
}

# So we can configure the inputs.
variable "role_id" {
  type = string
}
variable "role_name" {
  type = string
}

# So we can read output of the plan.
output "by_id" {
  value = data.idmc_role.by_id
}
output "by_name" {
  value = data.idmc_role.by_name
}
