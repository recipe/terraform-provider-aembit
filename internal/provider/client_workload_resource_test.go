package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClientWorkloadResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/TestAccClientWorkloadResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "unittest1namespace", fmt.Sprintf("unittest1namespace%d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "unittest1namespace", fmt.Sprintf("unittest1namespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr("aembit_client_workload.test", "name", "Unit Test 1"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "description", "Acceptance Test client workload"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr("aembit_client_workload.test", "identities.#", "1"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "identities.0.type", "k8sNamespace"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "identities.0.value", fmt.Sprintf("unittest1namespace%d", randID)),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_client_workload.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_client_workload.test", "id"),
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
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_client_workload.test", "name", "Unit Test 1 - modified"),
					// Verify Service Endpoint Host updated.
					resource.TestCheckResourceAttr("aembit_client_workload.test", "is_active", "false"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_client_workload.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_client_workload.test", "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
