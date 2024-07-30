# The correct provider source needs to be selected.
terraform {
  required_providers {
    idmc = {
      source = "tzrlk/idmc"
    }
  }
}

# So we can configure the inputs.
provider "idmc" {
}

# So we can read output of the plan.
output "example" {
  value = data.idmc_role_list.example
}
