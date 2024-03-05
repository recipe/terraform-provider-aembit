terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

resource "aembit_client_workload" "edu" {
  name        = "terraform client workload3"
  description = "new client workload3"
  is_active   = false
  identities = [
    {
      type  = "k8sNamespace"
      value = "workload3mod"
    },
  ]
}

output "edu_client_workload" {
  value = aembit_client_workload.edu
}

