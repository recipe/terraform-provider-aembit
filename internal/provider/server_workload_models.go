package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// serverWorkloadResourceModel maps the resource schema.
type serverWorkloadResourceModel struct {
	// ID is required for Framework acceptance testing
	ID              types.String          `tfsdk:"id"`
	Name            types.String          `tfsdk:"name"`
	Description     types.String          `tfsdk:"description"`
	IsActive        types.Bool            `tfsdk:"is_active"`
	ServiceEndpoint *serviceEndpointModel `tfsdk:"service_endpoint"`
}

// serverWorkloadDataSourceModel maps the datasource schema.
type serverWorkloadsDataSourceModel struct {
	ServerWorkloads []serverWorkloadResourceModel `tfsdk:"server_workloads"`
}

// serviceEndpointModel maps service endpoint data.
type serviceEndpointModel struct {
	ExternalID        types.String `tfsdk:"external_id"`
	ID                types.Int64  `tfsdk:"id"`
	Host              types.String `tfsdk:"host"`
	AppProtocol       types.String `tfsdk:"app_protocol"`
	TransportProtocol types.String `tfsdk:"transport_protocol"`
	RequestedPort     types.Int64  `tfsdk:"requested_port"`
	RequestedTLS      types.Bool   `tfsdk:"requested_tls"`
	Port              types.Int64  `tfsdk:"port"`
	TLS               types.Bool   `tfsdk:"tls"`

	WorkloadServiceAuthentication *workloadServiceAuthenticationModel `tfsdk:"workload_service_authentication"`
	TLSVerification               types.String                        `tfsdk:"tls_verification"`
}

// workloadServiceAuthenticationModel maps the WorkloadServiceAuthenticationDTO struct.
type workloadServiceAuthenticationModel struct {
	Method types.String `tfsdk:"method"`
	Scheme types.String `tfsdk:"scheme"`
	Config types.String `tfsdk:"config"`
}
