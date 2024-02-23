package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// agentControllerResourceModel maps the resource schema.
type agentControllerResourceModel struct {
	// ID is required for Framework acceptance testing
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	TrustProviderID types.String `tfsdk:"trust_provider_id"`
}

// agentControllerDataSourceModel maps the datasource schema.
type agentControllersDataSourceModel struct {
	AgentControllers []agentControllerResourceModel `tfsdk:"agent_controllers"`
}
