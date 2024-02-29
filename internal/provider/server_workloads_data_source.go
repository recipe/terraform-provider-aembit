package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverWorkloadsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverWorkloadsDataSource{}
)

// NewServerWorkloadsDataSource is a helper function to simplify the provider implementation.
func NewServerWorkloadsDataSource() datasource.DataSource {
	return &serverWorkloadsDataSource{}
}

// serverWorkloadsDataSource is the data source implementation.
type serverWorkloadsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *serverWorkloadsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.CloudClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *serverWorkloadsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_workloads"
}

// Schema defines the schema for the resource.
func (d *serverWorkloadsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an server workload.",
		Attributes: map[string]schema.Attribute{
			"server_workloads": schema.ListNestedAttribute{
				Description: "List of server workloads.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the server workload.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the server workload.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the server workload.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the server workload.",
							Computed:    true,
						},
						//"tags": schema.ListNestedAttribute{
						//	Description: "List of Tags.",
						//	Computed:    true,
						//	NestedObject: schema.NestedAttributeObject{
						//		Attributes: map[string]schema.Attribute{
						//			"key": schema.StringAttribute{
						//				Description: "Tag key.",
						//				Computed:    true,
						//			},
						//			"value": schema.StringAttribute{
						//				Description: "Tag value.",
						//				Computed:    true,
						//			},
						//		},
						//	},
						//},
						"service_endpoint": schema.SingleNestedAttribute{
							Description: "Service endpoint details.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Description: "Alphanumeric identifier of the service endpoint.",
									Computed:    true,
								},
								"id": schema.Int64Attribute{
									Description: "Number identifier of the service endpoint.",
									Computed:    true,
								},
								"host": schema.StringAttribute{
									Description: "hostname of the service endpoint.",
									Computed:    true,
								},
								"port": schema.Int64Attribute{
									Description: "port of the service endpoint.",
									Computed:    true,
								},
								"app_protocol": schema.StringAttribute{
									Description: "protocol of the service endpoint.",
									Computed:    true,
								},
								"requested_port": schema.Int64Attribute{
									Description: "requested port of the service endpoint.",
									Computed:    true,
								},
								"requested_tls": schema.BoolAttribute{
									Description: "requested tls of the service endpoint.",
									Computed:    true,
								},
								"tls_verification": schema.StringAttribute{
									Description: "tls verification of the service endpoint.",
									Computed:    true,
								},
								"transport_protocol": schema.StringAttribute{
									Description: "transport protocol of the service endpoint.",
									Computed:    true,
								},
								"tls": schema.BoolAttribute{
									Description: "tls of the service endpoint.",
									Computed:    true,
								},
								"workload_service_authentication": schema.SingleNestedAttribute{
									Description: "Service authentication details.",
									Computed:    true,
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"method": schema.StringAttribute{
											Description: "Service authentication method.",
											Computed:    true,
										},
										"scheme": schema.StringAttribute{
											Description: "Service authentication scheme.",
											Computed:    true,
										},
										"config": schema.StringAttribute{
											Description: "Service authentication config.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *serverWorkloadsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverWorkloadsDataSourceModel

	serverWorkloads, err := d.client.GetServerWorkloads(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Server Workloads",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, serverWorkload := range serverWorkloads {
		serverWorkloadState := serverWorkloadResourceModel{
			ID:          types.StringValue(serverWorkload.EntityDTO.ExternalID),
			Name:        types.StringValue(serverWorkload.EntityDTO.Name),
			Description: types.StringValue(serverWorkload.EntityDTO.Description),
			IsActive:    types.BoolValue(serverWorkload.EntityDTO.IsActive),
		}

		serverWorkloadState.ServiceEndpoint = &serviceEndpointModel{
			ExternalID:        types.StringValue(serverWorkload.ServiceEndpoint.ExternalID),
			ID:                types.Int64Value(int64(serverWorkload.ServiceEndpoint.ID)),
			Host:              types.StringValue(serverWorkload.ServiceEndpoint.Host),
			AppProtocol:       types.StringValue(serverWorkload.ServiceEndpoint.AppProtocol),
			TransportProtocol: types.StringValue(serverWorkload.ServiceEndpoint.TransportProtocol),
			RequestedPort:     types.Int64Value(int64(serverWorkload.ServiceEndpoint.RequestedPort)),
			RequestedTLS:      types.BoolValue(serverWorkload.ServiceEndpoint.RequestedTLS),
			Port:              types.Int64Value(int64(serverWorkload.ServiceEndpoint.Port)),
			TLS:               types.BoolValue(serverWorkload.ServiceEndpoint.TLS),
			TLSVerification:   types.StringValue(serverWorkload.ServiceEndpoint.TLSVerification),
		}

		state.ServerWorkloads = append(state.ServerWorkloads, serverWorkloadState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
