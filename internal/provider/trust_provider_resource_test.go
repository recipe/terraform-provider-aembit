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

func TestAccTrustProviderResource_AwsEcsRole(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/aws_ecs/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/aws_ecs/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.aws_ecs", "name", "TF Acceptance AWS ECS"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws_ecs", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws_ecs", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.aws_ecs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.aws_ecs", "name", "TF Acceptance AWS ECS - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws_ecs", "id"),
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

func TestAccTrustProviderResource_GcpIdentity(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/gcp/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/gcp/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.gcp", "name", "TF Acceptance GCP Identity"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.gcp", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.gcp", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.gcp",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.gcp", "name", "TF Acceptance GCP Identity - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.gcp", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_GitHubAction(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/github/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/github/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.github", "name", "TF Acceptance GitHub Action"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.github", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.github", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.github",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.github", "name", "TF Acceptance GitHub Action - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.github", "id"),
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
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "tags.day", "Sunday"),
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
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_KubernetesServiceAccount(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/kubernetes/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/kubernetes/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "name", "TF Acceptance Kubernetes"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes", "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "tags.color", "blue"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "tags.day", "Sunday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes_key", "name", "TF Acceptance Kubernetes Key"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.kubernetes",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "name", "TF Acceptance Kubernetes - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "tags.%", "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "tags.color", "orange"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "tags.day", "Tuesday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes_key", "name", "TF Acceptance Kubernetes Key - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_TerraformWorkspace(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/terraform/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/terraform/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.terraform", "name", "TF Acceptance Terraform Workspace"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.terraform", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.terraform", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_trust_provider.terraform",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.terraform", "name", "TF Acceptance Terraform Workspace - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.terraform", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
