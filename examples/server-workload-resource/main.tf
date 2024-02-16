terraform {
  required_providers {
    aembit = {
      source  = "aembit.io/dev/aembit"
    }
  }
}

provider "aembit" {
}

resource "aembit_server_workload" "edu" {
    name = "terraform server workload2"
    description = "new server workload"
    is_active = false
    service_endpoint = {
        host = "myhost.jason.com"
        port = 443
        app_protocol = "HTTP"
        transport_protocol = "TCP"
        requested_port = 80
        tls_verification = "full"
	    requested_tls = true
	    tls = true
    }
}

output "edu_server_workload" {
  value = aembit_server_workload.edu
}

