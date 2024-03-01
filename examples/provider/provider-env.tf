terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}


provider "aembit" {
}

resource "aembit_client_workload" "client" {
  # Resource configuration
}

