terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_access_conditions" "first" {}

output "first" {
  value = data.aembit_access_conditions.first
}
