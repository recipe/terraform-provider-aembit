package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// roleResourceModel maps the resource schema.
type roleResourceModel struct {
	// ID is required for Framework acceptance testing
	ID             types.String    `tfsdk:"id"`
	Name           types.String    `tfsdk:"name"`
	Description    types.String    `tfsdk:"description"`
	IsActive       types.Bool      `tfsdk:"is_active"`
	Tags           types.Map       `tfsdk:"tags"`
	AccessPolicies *rolePermission `tfsdk:"access_policies"`

	ClientWorkloads     *rolePermission `tfsdk:"client_workloads"`
	TrustProviders      *rolePermission `tfsdk:"trust_providers"`
	AccessConditions    *rolePermission `tfsdk:"access_conditions"`
	Integrations        *rolePermission `tfsdk:"integrations"`
	CredentialProviders *rolePermission `tfsdk:"credential_providers"`
	ServerWorkloads     *rolePermission `tfsdk:"server_workloads"`

	AgentControllers *rolePermission `tfsdk:"agent_controllers"`

	AccessAuthorizationEvents *roleReadOnlyPermission `tfsdk:"access_authorization_events"`
	AuditLogs                 *roleReadOnlyPermission `tfsdk:"audit_logs"`
	WorkloadEvents            *roleReadOnlyPermission `tfsdk:"workload_events"`

	Users             *rolePermission `tfsdk:"users"`
	Roles             *rolePermission `tfsdk:"roles"`
	LogStreams        *rolePermission `tfsdk:"log_streams"`
	IdentityProviders *rolePermission `tfsdk:"identity_providers"`
}

// roleDataSourceModel maps the datasource schema.
type rolesDataSourceModel struct {
	Roles []roleResourceModel `tfsdk:"roles"`
}

type rolePermission struct {
	Read  types.Bool `tfsdk:"read"`
	Write types.Bool `tfsdk:"write"`
}

type roleReadOnlyPermission struct {
	Read types.Bool `tfsdk:"read"`
}
