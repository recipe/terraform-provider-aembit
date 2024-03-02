terraform {
  required_providers {
    aembit = {
      source = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

resource "aembit_agent_controller" "azure_tp" {
  name      = "TF Acceptance Azure Trust Provider"
  is_active = true
  tags = {
    color = "blue"
    day   = "Sunday"
  }

  trust_provider_id = aembit_trust_provider.azure.id
}

resource "aembit_trust_provider" "azure" {
  name = "TF Acceptance Azure"
  azure_metadata = {
    subscription_id = "subscription_id"
  }
}

data "aembit_agent_controller_device_code" "first" {
  agent_controller_id = aembit_agent_controller.azure_tp.id
}

output "first" {
  value = data.aembit_agent_controller_device_code.first
}

