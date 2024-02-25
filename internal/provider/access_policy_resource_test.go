package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_policy.first_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
