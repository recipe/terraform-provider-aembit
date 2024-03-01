terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_integrations" "first" {}

output "first" {
  value     = data.aembit_integrations.first
  sensitive = true
}
