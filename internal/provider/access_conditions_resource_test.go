package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessConditionResource_Wiz(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/condition/wiz/TestAccAccessConditionResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/condition/wiz/TestAccAccessConditionResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr("aembit_access_condition.wiz", "name", "TF Acceptance Wiz"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_access_condition.wiz", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_access_condition.wiz", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_condition.wiz",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_access_condition.wiz", "name", "TF Acceptance Wiz - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAccessConditionResource_Crowdstrike(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/condition/crowdstrike/TestAccAccessConditionResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/condition/crowdstrike/TestAccAccessConditionResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr("aembit_access_condition.crowdstrike", "name", "TF Acceptance Crowdstrike"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_access_condition.crowdstrike", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_access_condition.crowdstrike", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_condition.crowdstrike",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_access_condition.crowdstrike", "name", "TF Acceptance Crowdstrike - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
