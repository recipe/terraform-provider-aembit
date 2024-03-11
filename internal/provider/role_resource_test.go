package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "TF Acceptance Role", fmt.Sprintf("TF Acceptance Role %d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "TF Acceptance Role", fmt.Sprintf("TF Acceptance Role %d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_role.role", "name", fmt.Sprintf("TF Acceptance Role %d", randID)),
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
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_role.role", "name", fmt.Sprintf("TF Acceptance Role %d - Modified", randID)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
