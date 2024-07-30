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

# So we can read output of the plan.
output "linux" {
  value = data.idmc_agent_installer.linux
}
output "windows" {
  value = data.idmc_agent_installer.windows
}
