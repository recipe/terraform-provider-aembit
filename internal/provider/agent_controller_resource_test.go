package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAgentControllerResource_WithTrustProvider(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/agent_controllers/TestAccAgentControllerResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/agent_controllers/TestAccAgentControllerResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "name", "TF Acceptance Azure Trust Provider"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_agent_controller.azure_tp", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_agent_controller.azure_tp", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_agent_controller.azure_tp",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "name", "TF Acceptance Azure Trust Provider - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_agent_controller.azure_tp", "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
