package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &agentControllerResource{}
	_ resource.ResourceWithConfigure   = &agentControllerResource{}
	_ resource.ResourceWithImportState = &agentControllerResource{}
)

// NewAgentControllerResource is a helper function to simplify the provider implementation.
func NewAgentControllerResource() resource.Resource {
	return &agentControllerResource{}
}

// agentControllerResource is the resource implementation.
type agentControllerResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *agentControllerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_controller"
}

// Configure adds the provider configured client to the resource.
func (r *agentControllerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *agentControllerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Agent Controller.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Agent Controller.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Agent Controller.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Agent Controller.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"trust_provider_id": schema.StringAttribute{
				Description: "Unique Trust Provider to use for authentication of the Agent Controller.",
				Optional:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *agentControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan agentControllerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var controller aembit.AgentControllerDTO = convertAgentControllerModelToDTO(ctx, plan, nil)

	// Create new Agent Controller
	agentController, err := r.client.CreateAgentController(controller, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Agent Controller",
			"Could not create Agent Controller, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertAgentControllerDTOToModel(ctx, *agentController)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *agentControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state agentControllerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed controller value from Aembit
	agentController, err := r.client.GetAgentController(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Agent Controller",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertAgentControllerDTOToModel(ctx, agentController)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *agentControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state agentControllerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan agentControllerResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var controller aembit.AgentControllerDTO = convertAgentControllerModelToDTO(ctx, plan, &externalID)

	// Update Agent Controller
	agentController, err := r.client.UpdateAgentController(controller, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Agent Controller",
			"Could not update Agent Controller, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertAgentControllerDTOToModel(ctx, *agentController)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *agentControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state agentControllerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Agent Controller is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableAgentController(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Agent Controller",
				"Could not disable Agent Controller, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Agent Controller
	_, err := r.client.DeleteAgentController(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Agent Controller",
			"Could not delete Agent Controller, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *agentControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertAgentControllerModelToDTO(ctx context.Context, model agentControllerResourceModel, externalID *string) aembit.AgentControllerDTO {
	var controller aembit.AgentControllerDTO
	controller.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			controller.Tags = append(controller.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}
	if externalID != nil {
		controller.EntityDTO.ExternalID = *externalID
	}
	controller.TrustProviderID = model.TrustProviderID.ValueString()

	return controller
}

func convertAgentControllerDTOToModel(ctx context.Context, dto aembit.AgentControllerDTO) agentControllerResourceModel {
	var model agentControllerResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	if len(dto.TrustProviderID) > 0 {
		model.TrustProviderID = types.StringValue(dto.TrustProviderID)
	} else {
		model.TrustProviderID = types.StringNull()
	}
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	return model
}
