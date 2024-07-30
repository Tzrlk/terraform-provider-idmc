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
