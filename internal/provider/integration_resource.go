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
	_ resource.Resource                = &integrationResource{}
	_ resource.ResourceWithConfigure   = &integrationResource{}
	_ resource.ResourceWithImportState = &integrationResource{}
)

// NewIntegrationResource is a helper function to simplify the provider implementation.
func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

// integrationResource is the resource implementation.
type integrationResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

// Configure adds the provider configured client to the resource.
func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the integration.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the integration.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the integration.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the integration.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of Aembit integration (either `WizIntegrationApi` or `CrowdStrike`).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"WizIntegrationApi", "CrowdStrike"}...),
				},
			},
			"sync_frequency": schema.Int64Attribute{
				Description: "Frequency to be used for synchronizing the integration.",
				Required:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "Endpoint to be used for performing the integration.",
				Required:    true,
			},
			"oauth_client_credentials": schema.SingleNestedAttribute{
				Description: "OAuth Client Credentials authentication information for the integration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"token_url": schema.StringAttribute{Required: true},
					"client_id": schema.StringAttribute{Required: true},
					"client_secret": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
					"audience": schema.StringAttribute{Optional: true},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan integrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.IntegrationDTO = convertIntegrationModelToDTO(ctx, plan, nil)

	// Create new Integration
	integration, err := r.client.CreateIntegration(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating integration",
			"Could not create integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = ConvertIntegrationDTOToModel(ctx, *integration, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state integrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	integration, err := r.client.GetIntegration(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Integration",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = ConvertIntegrationDTOToModel(ctx, integration, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state integrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan integrationResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.IntegrationDTO = convertIntegrationModelToDTO(ctx, plan, &externalID)

	// Update Integration
	integration, err := r.client.UpdateIntegration(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating integration",
			"Could not update integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = ConvertIntegrationDTOToModel(ctx, *integration, state)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state integrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Integration is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Integration",
			"Integration is active and cannot be deleted. Please mark the Integration as inactive first.",
		)
		return
	}

	// Delete existing Integration
	_, err := r.client.DeleteIntegration(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Integration",
			"Could not delete integration, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertIntegrationModelToDTO(ctx context.Context, model integrationResourceModel, externalID *string) aembit.IntegrationDTO {
	var integration aembit.IntegrationDTO
	integration.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			integration.Tags = append(integration.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}

	if externalID != nil {
		integration.EntityDTO.ExternalID = *externalID
	}

	integration.Endpoint = model.Endpoint.ValueString()
	integration.Type = model.Type.ValueString()
	integration.SyncFrequencySeconds = model.SyncFrequency.ValueInt64()
	integration.IntegrationJSON = aembit.IntegrationJSONDTO{
		TokenURL:     model.OAuthClientCredentials.TokenURL.ValueString(),
		ClientID:     model.OAuthClientCredentials.ClientID.ValueString(),
		ClientSecret: model.OAuthClientCredentials.ClientSecret.ValueString(),
		Audience:     model.OAuthClientCredentials.Audience.ValueString(),
	}

	return integration
}

func ConvertIntegrationDTOToModel(ctx context.Context, dto aembit.IntegrationDTO, state integrationResourceModel) integrationResourceModel {
	var model integrationResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	model.Type = types.StringValue(dto.Type)
	model.Endpoint = types.StringValue(dto.Endpoint)
	model.SyncFrequency = types.Int64Value(dto.SyncFrequencySeconds)
	model.OAuthClientCredentials = &integrationOAuthClientCredentialsModel{
		TokenURL:     types.StringValue(dto.IntegrationJSON.TokenURL),
		ClientID:     types.StringValue(dto.IntegrationJSON.ClientID),
		ClientSecret: types.StringValue(dto.IntegrationJSON.ClientSecret),
		Audience:     types.StringValue(dto.IntegrationJSON.Audience),
	}
	if len(dto.IntegrationJSON.ClientSecret) == 0 && state.OAuthClientCredentials != nil {
		model.OAuthClientCredentials.ClientSecret = state.OAuthClientCredentials.ClientSecret
	}

	return model
}
