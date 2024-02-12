package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// serverWorkloadResourceModel maps the resource schema.
type serverWorkloadResourceModel struct {
	ExternalId      types.String          `tfsdk:"external_id"`
	Name            types.String          `tfsdk:"name"`
	ServiceEndpoint *serviceEndpointModel `tfsdk:"service_endpoint"`
	Type            types.String          `tfsdk:"type"`

	// ID is required for Framework acceptance testing
	ID types.String `tfsdk:"id"`
}

// serverWorkloadDataSourceModel maps the datasource schema.
type serverWorkloadsDataSourceModel struct {
	ServerWorkloads []serverWorkloadResourceModel `tfsdk:"server_workloads"`
}

// serviceEndpointModel maps service endpoint data.
type serviceEndpointModel struct {
	ExternalId        types.String `tfsdk:"external_id"`
	Host              types.String `tfsdk:"host"`
	Port              types.Int64  `tfsdk:"port"`
	AppProtocol       types.String `tfsdk:"app_protocol"`
	TransportProtocol types.String `tfsdk:"transport_protocol"`
	RequestedPort     types.Int64  `tfsdk:"requested_port"`
	TlsVerification   types.String `tfsdk:"tls_verification"`
}
