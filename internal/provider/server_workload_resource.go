package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverWorkloadResource{}
	_ resource.ResourceWithConfigure   = &serverWorkloadResource{}
	_ resource.ResourceWithImportState = &serverWorkloadResource{}
)

// NewServerWorkloadResource is a helper function to simplify the provider implementation.
func NewServerWorkloadResource() resource.Resource {
	return &serverWorkloadResource{}
}

// serverWorkloadResource is the resource implementation.
type serverWorkloadResource struct {
	client *aembit.AembitClient
}

// Metadata returns the resource type name.
func (r *serverWorkloadResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_workload"
}

// Configure adds the provider configured client to the resource.
func (r *serverWorkloadResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.AembitClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.AembitClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *serverWorkloadResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the server workload.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the server workload.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the server workload.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the server workload.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of server workload.",
				Computed:    true,
			},
			"service_endpoint": schema.SingleNestedAttribute{
				Description: "Service endpoint details.",
				Required:    true,
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
						Required:    true,
					},
					"port": schema.Int64Attribute{
						Description: "port of the service endpoint.",
						Required:    true,
					},
					"app_protocol": schema.StringAttribute{
						Description: "protocol of the service endpoint.",
						Required:    true,
					},
					"requested_port": schema.Int64Attribute{
						Description: "requested port of the service endpoint.",
						Required:    true,
					},
					"tls_verification": schema.StringAttribute{
						Description: "tls verification of the service endpoint.",
						Required:    true,
					},
					"transport_protocol": schema.StringAttribute{
						Description: "transport protocol of the service endpoint.",
						Required:    true,
					},
					"requested_tls": schema.BoolAttribute{
						Description: "tls requested on the service endpoint.",
						Optional:    true,
						Computed:    true,
					},
					"tls": schema.BoolAttribute{
						Description: "tls indicated on the service endpoint.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *serverWorkloadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan serverWorkloadResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var workload aembit.ServerWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsActive:    plan.IsActive.ValueBool(),
	}

	workload.ServiceEndpoint = aembit.WorkloadServiceEndpointDTO{
		Host:              plan.ServiceEndpoint.Host.ValueString(),
		Id:                int(plan.ServiceEndpoint.Id.ValueInt64()),
		Port:              int(plan.ServiceEndpoint.Port.ValueInt64()),
		AppProtocol:       plan.ServiceEndpoint.AppProtocol.ValueString(),
		TransportProtocol: plan.ServiceEndpoint.TransportProtocol.ValueString(),
		RequestedPort:     int(plan.ServiceEndpoint.RequestedPort.ValueInt64()),
		RequestedTls:      plan.ServiceEndpoint.RequestedTls.ValueBool(),
		Tls:               plan.ServiceEndpoint.Tls.ValueBool(),
		TlsVerification:   plan.ServiceEndpoint.TlsVerification.ValueString(),
	}

	// Create new Server Workload
	server_workload, err := r.client.CreateServerWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server workload",
			"Could not create server workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(server_workload.EntityDTO.ExternalId)
	plan.Name = types.StringValue(server_workload.EntityDTO.Name)
	plan.Description = types.StringValue(server_workload.EntityDTO.Description)
	plan.IsActive = types.BoolValue(server_workload.EntityDTO.IsActive)
	plan.Type = types.StringValue(server_workload.Type)

	plan.ServiceEndpoint = &serviceEndpointModel{
		ExternalId:        types.StringValue(server_workload.ServiceEndpoint.ExternalId),
		Host:              types.StringValue(server_workload.ServiceEndpoint.Host),
		Port:              types.Int64Value(int64(server_workload.ServiceEndpoint.Port)),
		AppProtocol:       types.StringValue(server_workload.ServiceEndpoint.AppProtocol),
		TransportProtocol: types.StringValue(server_workload.ServiceEndpoint.TransportProtocol),
		RequestedPort:     types.Int64Value(int64(server_workload.ServiceEndpoint.RequestedPort)),
		RequestedTls:      types.BoolValue(server_workload.ServiceEndpoint.RequestedTls),
		Tls:               types.BoolValue(server_workload.ServiceEndpoint.Tls),
		TlsVerification:   types.StringValue(server_workload.ServiceEndpoint.TlsVerification),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serverWorkloadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state serverWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workload value from Aembit
	server_workload, err := r.client.GetServerWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Server Workload",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(server_workload.EntityDTO.ExternalId)
	state.Name = types.StringValue(server_workload.EntityDTO.Name)
	state.Description = types.StringValue(server_workload.EntityDTO.Description)
	state.IsActive = types.BoolValue(server_workload.EntityDTO.IsActive)
	state.Type = types.StringValue(server_workload.Type)

	state.ServiceEndpoint = &serviceEndpointModel{
		ExternalId:        types.StringValue(server_workload.ServiceEndpoint.ExternalId),
		Host:              types.StringValue(server_workload.ServiceEndpoint.Host),
		Port:              types.Int64Value(int64(server_workload.ServiceEndpoint.Port)),
		AppProtocol:       types.StringValue(server_workload.ServiceEndpoint.AppProtocol),
		TransportProtocol: types.StringValue(server_workload.ServiceEndpoint.TransportProtocol),
		RequestedPort:     types.Int64Value(int64(server_workload.ServiceEndpoint.RequestedPort)),
		RequestedTls:      types.BoolValue(server_workload.ServiceEndpoint.RequestedTls),
		Tls:               types.BoolValue(server_workload.ServiceEndpoint.Tls),
		TlsVerification:   types.StringValue(server_workload.ServiceEndpoint.TlsVerification),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverWorkloadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state serverWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	var external_id string
	external_id = state.ID.ValueString()

	// Retrieve values from plan
	var plan serverWorkloadResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var workload aembit.ServerWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		ExternalId:  external_id,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsActive:    plan.IsActive.ValueBool(),
	}
	workload.ServiceEndpoint = aembit.WorkloadServiceEndpointDTO{
		Host:              plan.ServiceEndpoint.Host.ValueString(),
		Port:              int(plan.ServiceEndpoint.Port.ValueInt64()),
		AppProtocol:       plan.ServiceEndpoint.AppProtocol.ValueString(),
		TransportProtocol: plan.ServiceEndpoint.TransportProtocol.ValueString(),
		RequestedPort:     int(plan.ServiceEndpoint.RequestedPort.ValueInt64()),
		RequestedTls:      plan.ServiceEndpoint.RequestedTls.ValueBool(),
		Tls:               plan.ServiceEndpoint.Tls.ValueBool(),
		TlsVerification:   plan.ServiceEndpoint.TlsVerification.ValueString(),
	}

	// Update Server Workload
	server_workload, err := r.client.UpdateServerWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating server workload",
			"Could not update server workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(server_workload.EntityDTO.ExternalId)
	plan.Name = types.StringValue(server_workload.EntityDTO.Name)
	plan.Description = types.StringValue(server_workload.EntityDTO.Description)
	plan.IsActive = types.BoolValue(server_workload.EntityDTO.IsActive)
	plan.Type = types.StringValue(server_workload.Type)
	plan.ServiceEndpoint = &serviceEndpointModel{
		ExternalId:        types.StringValue(server_workload.ServiceEndpoint.ExternalId),
		Host:              types.StringValue(server_workload.ServiceEndpoint.Host),
		Port:              types.Int64Value(int64(server_workload.ServiceEndpoint.Port)),
		AppProtocol:       types.StringValue(server_workload.ServiceEndpoint.AppProtocol),
		TransportProtocol: types.StringValue(server_workload.ServiceEndpoint.TransportProtocol),
		RequestedPort:     types.Int64Value(int64(server_workload.ServiceEndpoint.RequestedPort)),
		TlsVerification:   types.StringValue(server_workload.ServiceEndpoint.TlsVerification),
		RequestedTls:      types.BoolValue(server_workload.ServiceEndpoint.RequestedTls),
		Tls:               types.BoolValue(server_workload.ServiceEndpoint.Tls),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serverWorkloadResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state serverWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Server Workload is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Server Workload",
			"Server Workload is active and cannot be deleted. Please mark the workload as inactive first.",
		)
		return
	}

	// Delete existing Server Workload
	_, err := r.client.DeleteServerWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Server Workload",
			"Could not delete server workload, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId
func (r *serverWorkloadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
