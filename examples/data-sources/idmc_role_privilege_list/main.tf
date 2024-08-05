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
