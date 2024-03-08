package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAgentControllerDeviceCodeDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/device_code/TestAccDeviceCodeDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_agent_controller_device_code.test", "device_code"),
				),
			},
		},
	})
}
