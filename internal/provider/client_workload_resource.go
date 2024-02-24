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
	client *aembit.CloudClient
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
			"identities": schema.SetNestedAttribute{
				Description: "Set of client workload identities.",
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
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
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
	var workload aembit.ClientWorkloadExternalDTO = convertClientWorkloadModelToDTO(ctx, plan, nil)

	// Create new Client Workload
	clientWorkload, err := r.client.CreateClientWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating client workload",
			"Could not create client workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertClientWorkloadDTOToModel(ctx, *clientWorkload)

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
	clientWorkload, err := r.client.GetClientWorkload(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Client Workload",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = convertClientWorkloadDTOToModel(ctx, clientWorkload)

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
	var externalID string = state.ID.ValueString()

	// Retrieve values from plan
	var plan clientWorkloadResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var workload aembit.ClientWorkloadExternalDTO = convertClientWorkloadModelToDTO(ctx, plan, &externalID)

	// Update Client Workload
	clientWorkload, err := r.client.UpdateClientWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating client workload",
			"Could not update client workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertClientWorkloadDTOToModel(ctx, *clientWorkload)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
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

// Imports an existing resource by passing externalId.
func (r *clientWorkloadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertClientWorkloadModelToDTO(ctx context.Context, model clientWorkloadResourceModel, externalID *string) aembit.ClientWorkloadExternalDTO {
	var workload aembit.ClientWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	var identities []identitiesModel
	if len(model.Identities.Elements()) > 0 {
		_ = model.Identities.ElementsAs(ctx, &identities, false)

		for _, identity := range identities {
			workload.Identities = append(workload.Identities, aembit.ClientWorkloadIdentityDTO{
				Type:  identity.Type.ValueString(),
				Value: identity.Value.ValueString(),
			})
		}

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

	if externalID != nil {
		workload.EntityDTO.ExternalID = *externalID
	}

	return workload
}

func convertClientWorkloadDTOToModel(ctx context.Context, dto aembit.ClientWorkloadExternalDTO) clientWorkloadResourceModel {
	var model clientWorkloadResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Type = types.StringValue(dto.Type)
	model.Identities = newClientWorkloadIdentityModel(ctx, dto.Identities)
	if len(dto.EntityDTO.Tags) > 0 {
		model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)
	}

	return model
}

func newClientWorkloadIdentityModel(ctx context.Context, clientWorkloadIdentities []aembit.ClientWorkloadIdentityDTO) types.Set {
	identities := make([]identitiesModel, len(clientWorkloadIdentities))

	for i, identity := range clientWorkloadIdentities {
		identities[i] = identitiesModel{
			Type:  types.StringValue(identity.Type),
			Value: types.StringValue(identity.Value),
		}
	}

	s, _ := types.SetValueFrom(ctx, TfIdentityObjectType, identities)
	return s
}
