package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClientWorkloadResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/TestAccClientWorkloadResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr("aembit_client_workload.test", "name", "Unit Test 1"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "description", "Acceptance Test client workload"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "is_active", "true"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr("aembit_client_workload.test", "identities.#", "1"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "identities.0.type", "k8sNamespace"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "identities.0.value", "unittest1namespace"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_client_workload.test", "id"),
					resource.TestCheckResourceAttrSet("aembit_client_workload.test", "type"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_client_workload.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_client_workload.test", "name", "Unit Test 1 - modified"),
					// Verify Service Endpoint Host updated.
					resource.TestCheckResourceAttr("aembit_client_workload.test", "is_active", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
