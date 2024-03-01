terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_agent_controllers" "first" {}

output "first" {
  value = data.aembit_agent_controllers.first
}
