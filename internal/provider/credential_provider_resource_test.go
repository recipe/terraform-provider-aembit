package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialProviderResource_ApiKey(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/apikey/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/apikey/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.api_key", "name", "TF Acceptance API Key"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.api_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.api_key", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.api_key",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.api_key", "name", "TF Acceptance API Key - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_OAuthClientCredentials(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/oauth/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth", "name", "TF Acceptance OAuth"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.oauth", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.oauth",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.oauth", "name", "TF Acceptance OAuth - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCredentialProviderResource_VaultClientToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/credential/vault/TestAccCredentialProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/credential/vault/TestAccCredentialProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Credential Provider Name
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "name", "TF Acceptance Vault"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_credential_provider.vault", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_credential_provider.vault", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_credential_provider.vault",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_credential_provider.vault", "name", "TF Acceptance Vault - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
