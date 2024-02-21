package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustProviderResource_AzureMetadata(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/azure/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/azure/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.azure", "name", "TF Acceptance Azure"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.azure", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.azure", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.azure",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.azure", "name", "TF Acceptance Azure - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_AwsMetadata(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/aws/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/aws/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.aws", "name", "TF Acceptance AWS"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.aws",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.aws", "name", "TF Acceptance AWS - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_Kerberos(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/kerberos/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/kerberos/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "name", "TF Acceptance Kerberos"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kerberos", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kerberos", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.kerberos",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "name", "TF Acceptance Kerberos - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
