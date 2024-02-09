package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerWorkloadResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "aembit_server_workload" "test" {
	name = "Unit Test 1"
	service_endpoint = {
		host = "unittest.testhost.com"
		port = 443
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 80
		tls_verification = "full"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder ID
					resource.TestCheckResourceAttr("aembit_server_workload.test", "id", "placeholder"),
					// Verify Server Workload Name
					resource.TestCheckResourceAttr("aembit_server_workload.test", "name", "Unit Test 1"),
					// Verify Service Endpoint.
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.host", "unittest.testhost.com"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.port", "443"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.app_protocol", "HTTP"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.transport_protocol", "TCP"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.requested_port", "80"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.tls_verification", "full"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "external_id"),
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "type"),
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "service_endpoint.external_id"),
				),
			},
			// ImportState testing
			//{
			//	ResourceName:      "aembit_server_workload.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "aembit_server_workload" "test" {
	name = "Unit Test 1 - Modified"
	service_endpoint = {
		host = "unittest.testhost2.com"
		port = 443
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 80
		tls_verification = "full"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_server_workload.test", "name", "Unit Test 1 - Modified"),
					// Verify Service Endpoint Host updated.
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.host", "unittest.testhost2.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
