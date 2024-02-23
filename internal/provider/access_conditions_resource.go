package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &accessConditionResource{}
	_ resource.ResourceWithConfigure   = &accessConditionResource{}
	_ resource.ResourceWithImportState = &accessConditionResource{}
)

// NewAccessConditionResource is a helper function to simplify the provider implementation.
func NewAccessConditionResource() resource.Resource {
	return &accessConditionResource{}
}

// accessConditionResource is the resource implementation.
type accessConditionResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *accessConditionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_condition"
}

// Configure adds the provider configured client to the resource.
func (r *accessConditionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *accessConditionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the Access Condition.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the Access Condition.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the Access Condition.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the Access Condition.",
				Optional:    true,
				Computed:    true,
			},
			"integration_id": schema.StringAttribute{
				Description: "ID of the Integration used by the Access Condition.",
				Required:    true,
			},
			"wiz_conditions": schema.SingleNestedAttribute{
				Description: "Wiz Specific rules for the Access Condition.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"max_last_seen":               schema.Int64Attribute{Required: true},
					"container_cluster_connected": schema.BoolAttribute{Required: true},
				},
			},
			"crowdstrike_conditions": schema.SingleNestedAttribute{
				Description: "CrowdStrike Specific rules for the Access Condition.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"max_last_seen":       schema.Int64Attribute{Required: true},
					"match_hostname":      schema.BoolAttribute{Required: true},
					"match_serial_number": schema.BoolAttribute{Required: true},
					"prevent_rfm":         schema.BoolAttribute{Required: true},
				},
			},
		},
	}
}

// Configure validators to ensure that only one trust provider type is specified.
func (r *accessConditionResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("wiz_conditions"),
			path.MatchRoot("crowdstrike_conditions"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *accessConditionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accessConditionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.AccessConditionDTO = convertAccessConditionModelToDTO(ctx, plan, nil)

	// Create new AccessCondition
	accessCondition, err := r.client.CreateAccessCondition(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Access Condition",
			"Could not create Access Condition, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertAccessConditionDTOToModel(ctx, *accessCondition, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *accessConditionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state accessConditionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	accessCondition, err := r.client.GetAccessCondition(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Access Condition",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertAccessConditionDTOToModel(ctx, accessCondition, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *accessConditionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state accessConditionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan accessConditionResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.AccessConditionDTO = convertAccessConditionModelToDTO(ctx, plan, &externalID)

	// Update AccessCondition
	accessCondition, err := r.client.UpdateAccessCondition(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Access Condition",
			"Could not update Access Condition, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertAccessConditionDTOToModel(ctx, *accessCondition, state)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *accessConditionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state accessConditionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if AccessCondition is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Access Condition",
			"Access Condition is active and cannot be deleted. Please mark the Access Condition as inactive first.",
		)
		return
	}

	// Delete existing AccessCondition
	_, err := r.client.DeleteAccessCondition(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AccessCondition",
			"Could not delete Access Condition, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *accessConditionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertAccessConditionModelToDTO(_ context.Context, model accessConditionResourceModel, externalID *string) aembit.AccessConditionDTO {
	var accessCondition aembit.AccessConditionDTO
	accessCondition.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if externalID != nil {
		accessCondition.EntityDTO.ExternalID = *externalID
	}

	accessCondition.IntegrationID = model.IntegrationID.ValueString()
	if model.Wiz != nil {
		accessCondition.Conditions.MaxLastSeenSeconds = model.Wiz.MaxLastSeen.ValueInt64()
		accessCondition.Conditions.ContainerClusterConnected = model.Wiz.ContainerClusterConnected.ValueBool()
	}
	if model.CrowdStrike != nil {
		accessCondition.Conditions.MaxLastSeenSeconds = model.CrowdStrike.MaxLastSeen.ValueInt64()
		accessCondition.Conditions.MatchHostname = model.CrowdStrike.MatchHostname.ValueBool()
		accessCondition.Conditions.MatchSerialNumber = model.CrowdStrike.MatchSerialNumber.ValueBool()
		accessCondition.Conditions.PreventRestrictedFunctionalityMode = model.CrowdStrike.PreventRestrictedFunctionalityMode.ValueBool()
	}

	return accessCondition
}

func convertAccessConditionDTOToModel(_ context.Context, dto aembit.AccessConditionDTO, _ accessConditionResourceModel) accessConditionResourceModel {
	var model accessConditionResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	if len(dto.IntegrationID) == 0 {
		model.IntegrationID = types.StringValue(dto.Integration.ExternalID)
	} else {
		model.IntegrationID = types.StringValue(dto.IntegrationID)
	}
	switch dto.Integration.Type {
	case "WizIntegrationApi":
		model.Wiz = &accessConditionWizModel{
			MaxLastSeen:               types.Int64Value(dto.Conditions.MaxLastSeenSeconds),
			ContainerClusterConnected: types.BoolValue(dto.Conditions.ContainerClusterConnected),
		}
	case "CrowdStrike":
		model.CrowdStrike = &accessConditionCrowdstrikeModel{
			MaxLastSeen:                        types.Int64Value(dto.Conditions.MaxLastSeenSeconds),
			MatchHostname:                      types.BoolValue(dto.Conditions.MatchHostname),
			MatchSerialNumber:                  types.BoolValue(dto.Conditions.MatchSerialNumber),
			PreventRestrictedFunctionalityMode: types.BoolValue(dto.Conditions.PreventRestrictedFunctionalityMode),
		}
	}

	return model
}
