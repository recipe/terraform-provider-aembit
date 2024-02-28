package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerWorkloadResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/server/TestAccServerWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/server/TestAccServerWorkloadResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Server Workload Name
					resource.TestCheckResourceAttr("aembit_server_workload.test", "name", "Unit Test 1"),
					// Verify Service Endpoint.
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.host", "unittest.testhost.com"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.port", "443"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.app_protocol", "HTTP"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.transport_protocol", "TCP"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.requested_port", "80"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.tls_verification", "full"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.workload_service_authentication.method", "HTTP Authentication"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.workload_service_authentication.scheme", "Bearer"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "id"),
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "type"),
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "service_endpoint.external_id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_server_workload.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_server_workload.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_server_workload.test", "name", "Unit Test 1 - Modified"),
					// Verify Service Endpoint updated.
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.host", "unittest.testhost2.com"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.workload_service_authentication.method", "HTTP Authentication"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.workload_service_authentication.scheme", "Header"),
					resource.TestCheckResourceAttr("aembit_server_workload.test", "service_endpoint.workload_service_authentication.config", "X-Vault-Token"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
