terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_trust_providers" "first" {}

output "first" {
  value = data.aembit_trust_providers.first
}
