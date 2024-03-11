package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_role.role", "name", "TF Acceptance Role"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_role.role", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_role.role", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_role.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_role.role", "name", "TF Acceptance Role - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
