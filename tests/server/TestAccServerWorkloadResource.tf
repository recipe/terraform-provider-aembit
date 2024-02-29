provider "aembit" {
}

resource "aembit_server_workload" "test" {
	name = "Unit Test 1"
	service_endpoint = {
		host = "unittest.testhost.com"
		port = 443
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 80
		tls_verification = "full"
		authentication_config = {
			"method" = "HTTP Authentication"
			"scheme" = "Bearer"
		}
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

