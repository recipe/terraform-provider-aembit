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
	client *aembit.CloudClient
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

	client, ok := req.ProviderData.(*aembit.CloudClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
				Description: "Unique identifier of the Server Workload.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Server Workload.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for the Server Workload.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Server Workload.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"service_endpoint": schema.SingleNestedAttribute{
				Description: "Service endpoint details.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Description: "Unique identifier of the service endpoint.",
						Computed:    true,
					},
					"id": schema.Int64Attribute{
						Description: "Number identifier of the service endpoint.",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "Hostname of the Server Workload service endpoint.",
						Required:    true,
					},
					"port": schema.Int64Attribute{
						Description: "Port of the Server Workload service endpoint.",
						Required:    true,
					},
					"app_protocol": schema.StringAttribute{
						Description: "Application Protocol of the Server Workload service endpoint. Possible values are: \n" +
							"\t* `Amazon Redshift`\n" +
							"\t* `HTTP`\n" +
							"\t* `MySQL`\n" +
							"\t* `PostgreSQL`\n" +
							"\t* `Redis`\n" +
							"\t* `Snowflake`\n",
						Required: true,
					},
					"requested_port": schema.Int64Attribute{
						Description: "Requested port of the Server Workload service endpoint.",
						Required:    true,
					},
					"tls_verification": schema.StringAttribute{
						Description: "TLS verification configuration of the Server Workload service endpoint. Possible values are `full` (default) or `none`.",
						Required:    true,
					},
					"transport_protocol": schema.StringAttribute{
						Description: "Transport protocol of the Server Workload service endpoint. This value must be set to the default `TCP`.",
						Required:    true,
					},
					"requested_tls": schema.BoolAttribute{
						Description: "TLS requested on the Server Workload service endpoint.",
						Optional:    true,
						Computed:    true,
					},
					"tls": schema.BoolAttribute{
						Description: "TLS indicated on the Server Workload service endpoint.",
						Optional:    true,
						Computed:    true,
					},
					"authentication_config": schema.SingleNestedAttribute{
						Description: "Service authentication details.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"method": schema.StringAttribute{
								Description: "Server Workload Service authentication method. Possible values are: \n" +
									"\t* `API Key`\n" +
									"\t* `HTTP Authentication`\n" +
									"\t* `JWT Token Authentication`\n" +
									"\t* `Password Authentication`\n",
								Required: true,
							},
							"scheme": schema.StringAttribute{
								Description: "Server Workload Service authentication scheme. Possible values are: \n" +
									"\t* For Authentation Method `API Key`:\n" +
									"\t\t* `Header`\n" +
									"\t\t* `Query Parameter`\n" +
									"\t* For Authentation Method `HTTP Authentication`:\n" +
									"\t\t* `Basic`\n" +
									"\t\t* `Bearer`\n" +
									"\t\t* `Header`\n" +
									"\t\t* `AWS Signature v4`\n" +
									"\t* For Authentation Method `JWT Token Authentication`:\n" +
									"\t\t* `Snowflake JWT`\n" +
									"\t* For Authentation Method `Password Authentication`:\n" +
									"\t\t* `Password`\n",
								Required: true,
							},
							"config": schema.StringAttribute{
								Description: "Server Workload Service authentication config. This value is used to identify the HTTP Header or Query Parameter used for the associated authentication scheme.",
								Optional:    true,
								Computed:    true,
							},
						},
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
	var workload aembit.ServerWorkloadExternalDTO = convertServerWorkloadModelToDTO(ctx, plan, nil)

	// Create new Server Workload
	serverWorkload, err := r.client.CreateServerWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server workload",
			"Could not create server workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertServerWorkloadDTOToModel(ctx, *serverWorkload)

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
	serverWorkload, err := r.client.GetServerWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Server Workload",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = convertServerWorkloadDTOToModel(ctx, serverWorkload)

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
	var externalID string = state.ID.ValueString()

	// Retrieve values from plan
	var plan serverWorkloadResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var workload aembit.ServerWorkloadExternalDTO = convertServerWorkloadModelToDTO(ctx, plan, &externalID)

	// Update Server Workload
	serverWorkload, err := r.client.UpdateServerWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating server workload",
			"Could not update server workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertServerWorkloadDTOToModel(ctx, *serverWorkload)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
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

// Imports an existing resource by passing externalID.
func (r *serverWorkloadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertServerWorkloadModelToDTO(ctx context.Context, model serverWorkloadResourceModel, externalID *string) aembit.ServerWorkloadExternalDTO {
	var workload aembit.ServerWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			workload.Tags = append(workload.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}

	workload.ServiceEndpoint = aembit.WorkloadServiceEndpointDTO{
		Host:              model.ServiceEndpoint.Host.ValueString(),
		ID:                int(model.ServiceEndpoint.ID.ValueInt64()),
		Port:              int(model.ServiceEndpoint.Port.ValueInt64()),
		AppProtocol:       model.ServiceEndpoint.AppProtocol.ValueString(),
		TransportProtocol: model.ServiceEndpoint.TransportProtocol.ValueString(),
		RequestedPort:     int(model.ServiceEndpoint.RequestedPort.ValueInt64()),
		RequestedTLS:      model.ServiceEndpoint.RequestedTLS.ValueBool(),
		TLS:               model.ServiceEndpoint.TLS.ValueBool(),
		TLSVerification:   model.ServiceEndpoint.TLSVerification.ValueString(),
	}

	if model.ServiceEndpoint.WorkloadServiceAuthentication != nil {
		workload.ServiceEndpoint.WorkloadServiceAuthentication = &aembit.WorkloadServiceAuthenticationDTO{
			Method: model.ServiceEndpoint.WorkloadServiceAuthentication.Method.ValueString(),
			Scheme: model.ServiceEndpoint.WorkloadServiceAuthentication.Scheme.ValueString(),
			Config: model.ServiceEndpoint.WorkloadServiceAuthentication.Config.ValueString(),
		}
	}

	if externalID != nil {
		workload.EntityDTO.ExternalID = *externalID
	}

	return workload
}

func convertServerWorkloadDTOToModel(ctx context.Context, dto aembit.ServerWorkloadExternalDTO) serverWorkloadResourceModel {
	var model serverWorkloadResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	model.ServiceEndpoint = &serviceEndpointModel{
		ExternalID:        types.StringValue(dto.ServiceEndpoint.ExternalID),
		Host:              types.StringValue(dto.ServiceEndpoint.Host),
		Port:              types.Int64Value(int64(dto.ServiceEndpoint.Port)),
		AppProtocol:       types.StringValue(dto.ServiceEndpoint.AppProtocol),
		TransportProtocol: types.StringValue(dto.ServiceEndpoint.TransportProtocol),
		RequestedPort:     types.Int64Value(int64(dto.ServiceEndpoint.RequestedPort)),
		RequestedTLS:      types.BoolValue(dto.ServiceEndpoint.RequestedTLS),
		TLS:               types.BoolValue(dto.ServiceEndpoint.TLS),
		TLSVerification:   types.StringValue(dto.ServiceEndpoint.TLSVerification),
	}

	if dto.ServiceEndpoint.WorkloadServiceAuthentication != nil {
		model.ServiceEndpoint.WorkloadServiceAuthentication = &workloadServiceAuthenticationModel{
			Scheme: types.StringValue(dto.ServiceEndpoint.WorkloadServiceAuthentication.Scheme),
			Method: types.StringValue(dto.ServiceEndpoint.WorkloadServiceAuthentication.Method),
			Config: types.StringValue(dto.ServiceEndpoint.WorkloadServiceAuthentication.Config),
		}
	}

	return model
}
