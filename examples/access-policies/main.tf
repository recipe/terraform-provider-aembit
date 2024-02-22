terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_access_policies" "first" {}

output "first_access_policies" {
  value = data.aembit_access_policies.first
}
