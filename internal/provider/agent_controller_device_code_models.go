package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// agentControllerDeviceCodeDataSourceModel maps the resource schema.
type agentControllerDeviceCodeDataSourceModel struct {
	ID         types.String `tfsdk:"agent_controller_id"`
	DeviceCode types.String `tfsdk:"device_code"`
}
