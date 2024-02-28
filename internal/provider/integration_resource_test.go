package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationResource_Wiz(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/integration/wiz/TestAccIntegrationResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/integration/wiz/TestAccIntegrationResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr("aembit_integration.wiz", "name", "TF Acceptance Wiz"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_integration.wiz", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_integration.wiz", "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_integration.wiz", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_integration.wiz", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_integration.wiz", "tags.day", "Sunday"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_integration.wiz",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_integration.wiz", "name", "TF Acceptance Wiz - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_integration.wiz", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_integration.wiz", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_integration.wiz", "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationResource_Crowdstrike(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/integration/crowdstrike/TestAccIntegrationResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/integration/crowdstrike/TestAccIntegrationResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr("aembit_integration.crowdstrike", "name", "TF Acceptance Crowdstrike"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_integration.crowdstrike", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_integration.crowdstrike", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_integration.crowdstrike",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_integration.crowdstrike", "name", "TF Acceptance Crowdstrike - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
