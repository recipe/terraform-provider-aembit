package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clientWorkloadResource{}
	_ resource.ResourceWithConfigure   = &clientWorkloadResource{}
	_ resource.ResourceWithImportState = &clientWorkloadResource{}
)

// NewClientWorkloadResource is a helper function to simplify the provider implementation.
func NewClientWorkloadResource() resource.Resource {
	return &clientWorkloadResource{}
}

// clientWorkloadResource is the resource implementation.
type clientWorkloadResource struct {
	client *aembit.Client
}

// Metadata returns the resource type name.
func (r *clientWorkloadResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_workload"
}

// Configure adds the provider configured client to the resource.
func (r *clientWorkloadResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *clientWorkloadResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the client workload.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the client workload.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the client workload.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the client workload.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of client workload.",
				Computed:    true,
			},
			"identities": schema.ListNestedAttribute{
				Description: "List of client workload identities.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Client identity type.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "Client identity value.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *clientWorkloadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan clientWorkloadResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var workload aembit.ClientWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsActive:    plan.IsActive.ValueBool(),
	}

	for _, identity := range plan.Identities {
		workload.Identities = append(workload.Identities, aembit.ClientWorkloadIdentityDTO{
			Type:  identity.Type.ValueString(),
			Value: identity.Value.ValueString(),
		})
	}

	// Create new Client Workload
	client_workload, err := r.client.CreateClientWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating client workload",
			"Could not create client workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(client_workload.EntityDTO.ExternalId)
	plan.Name = types.StringValue(client_workload.EntityDTO.Name)
	plan.Description = types.StringValue(client_workload.EntityDTO.Description)
	plan.IsActive = types.BoolValue(client_workload.EntityDTO.IsActive)
	plan.Type = types.StringValue(client_workload.Type)

	for identityIndex, identityItem := range client_workload.Identities {
		plan.Identities[identityIndex] = identitiesModel{
			Type:  types.StringValue(identityItem.Type),
			Value: types.StringValue(identityItem.Value),
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clientWorkloadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state clientWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workload value from Aembit
	client_workload, err := r.client.GetClientWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Client Workload",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(client_workload.EntityDTO.ExternalId)
	state.Name = types.StringValue(client_workload.EntityDTO.Name)
	state.Description = types.StringValue(client_workload.EntityDTO.Description)
	state.IsActive = types.BoolValue(client_workload.EntityDTO.IsActive)
	state.Type = types.StringValue(client_workload.Type)

	// Check for changes in the Identities list.
	var mismatch bool = false
	if len(state.Identities) != len(client_workload.Identities) {
		// The count of Identities on the backend is different from the Terraform state.
		// Replace the state with the backend configuration.
		tflog.Debug(ctx, "Count of Client Workload identities differs from Terraform state.")
		// Mark mismatch to override Terraform state with backend configuration.
		mismatch = true
	}
	// Compare the identity list contents between backend and Terraform state.
	// Only perform this check if we haven't already found a difference.
	if mismatch == false {
		identity_exists := make(map[string]string)
		for _, identity := range state.Identities {
			// Build map of identities from Terraform state.
			identity_exists[identity.Type.ValueString()] = identity.Value.ValueString()
		}
		for _, identityItem := range client_workload.Identities {
			// Compare retrieved identities with Terraform state.
			val, ok := identity_exists[identityItem.Type]
			if !ok || val != identityItem.Value {
				// Mismatch found.
				tflog.Debug(ctx, "Client Workload identity list differs from Terraform state.")
				mismatch = true
				break
			}
		}
	}
	if mismatch == true {
		// Backend doesn't match Terraform state--replace Terraform state.
		var newIdentities []identitiesModel
		for _, identityItem := range client_workload.Identities {
			newIdentities = append(newIdentities, identitiesModel{
				Type:  types.StringValue(identityItem.Type),
				Value: types.StringValue(identityItem.Value),
			})
		}
		state.Identities = newIdentities
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clientWorkloadResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state clientWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	var external_id string
	external_id = state.ID.ValueString()

	// Retrieve values from plan
	var plan clientWorkloadResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var workload aembit.ClientWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		ExternalId:  external_id,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsActive:    plan.IsActive.ValueBool(),
	}

	for _, identity := range plan.Identities {
		workload.Identities = append(workload.Identities, aembit.ClientWorkloadIdentityDTO{
			Type:  identity.Type.ValueString(),
			Value: identity.Value.ValueString(),
		})
	}

	// Update Client Workload
	client_workload, err := r.client.UpdateClientWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating client workload",
			"Could not update client workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(client_workload.EntityDTO.ExternalId)
	plan.Name = types.StringValue(client_workload.EntityDTO.Name)
	plan.Description = types.StringValue(client_workload.EntityDTO.Description)
	plan.IsActive = types.BoolValue(client_workload.EntityDTO.IsActive)
	plan.Type = types.StringValue(client_workload.Type)

	for identityIndex, identityItem := range client_workload.Identities {
		plan.Identities[identityIndex] = identitiesModel{
			Type:  types.StringValue(identityItem.Type),
			Value: types.StringValue(identityItem.Value),
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clientWorkloadResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state clientWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Client Workload is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Client Workload",
			"Client Workload is active and cannot be deleted. Please mark the workload as inactive first.",
		)
		return
	}

	// Delete existing Client Workload
	_, err := r.client.DeleteClientWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Client Workload",
			"Could not delete client workload, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId
func (r *clientWorkloadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
