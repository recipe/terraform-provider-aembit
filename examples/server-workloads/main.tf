terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_server_workloads" "first" {}

output "first_server_workloads" {
  value = data.aembit_server_workloads.first
}
