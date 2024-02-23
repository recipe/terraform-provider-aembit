terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_client_workloads" "first" {}

output "first_client_workloads" {
  value = data.aembit_client_workloads.first
}
