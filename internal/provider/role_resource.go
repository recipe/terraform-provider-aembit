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
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

// NewRoleResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &roleResource{}
}

// roleResource is the resource implementation.
type roleResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Configure adds the provider configured client to the resource.
func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Role.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Role.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Role.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Role.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"access_policies":             definePermissionAttribute("Access Policy"),
			"client_workloads":            definePermissionAttribute("Client Workload"),
			"trust_providers":             definePermissionAttribute("Trust Provider"),
			"access_conditions":           definePermissionAttribute("Access Condition"),
			"integrations":                definePermissionAttribute("Integration"),
			"credential_providers":        definePermissionAttribute("Credential Provider"),
			"server_workloads":            definePermissionAttribute("Server Workload"),
			"agent_controllers":           definePermissionAttribute("Agent Controller"),
			"access_authorization_events": definePermissionReadOnlyAttribute("Access Authorization Event"),
			"audit_logs":                  definePermissionReadOnlyAttribute("Audit Log"),
			"workload_events":             definePermissionReadOnlyAttribute("Workload Event"),
			"users":                       definePermissionAttribute("User"),
			"roles":                       definePermissionAttribute("Role"),
			"log_streams":                 definePermissionAttribute("Log Stream"),
			"identity_providers":          definePermissionAttribute("Identity Provider"),
		},
	}
}

func definePermissionAttribute(name string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: fmt.Sprintf("Permissions for %s resources.", name),
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"read": schema.BoolAttribute{
				Description: fmt.Sprintf("Trust if this Role should be able to query and view %s resources.", name),
				Required:    true,
			},
			"write": schema.BoolAttribute{
				Description: fmt.Sprintf("True if this Role should be able to create and update %s resources.", name),
				Required:    true,
			},
		},
	}
}

func definePermissionReadOnlyAttribute(name string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: fmt.Sprintf("Permissions for %s resources.", name),
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"read": schema.BoolAttribute{
				Description: fmt.Sprintf("Trust if this Role should be able to query and view %s data.", name),
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan roleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.RoleDTO = convertRoleModelToDTO(ctx, plan, nil)

	// Create new Role
	role, err := r.client.CreateRole(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Role",
			"Could not create Role, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertRoleDTOToModel(ctx, *role)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state roleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	role, err := r.client.GetRole(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Role",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertRoleDTOToModel(ctx, role)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state roleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan roleResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.RoleDTO = convertRoleModelToDTO(ctx, plan, &externalID)

	// Update Role
	role, err := r.client.UpdateRole(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Role",
			"Could not update Role, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertRoleDTOToModel(ctx, *role)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state roleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Role is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableRole(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Role",
				"Could not disable Role, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Role
	_, err := r.client.DeleteRole(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Role",
			"Could not delete Role, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Model to DTO conversion methods.
func convertRoleModelToDTO(ctx context.Context, model roleResourceModel, externalID *string) aembit.RoleDTO {
	var dto aembit.RoleDTO
	dto.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			dto.Tags = append(dto.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}
	if externalID != nil {
		dto.EntityDTO.ExternalID = *externalID
	}

	dto.Permissions = make([]aembit.RolePermissionDTO, 0)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Access Policies", model.AccessPolicies)

	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Client Workloads", model.AccessPolicies)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Trust Providers", model.ClientWorkloads)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Access Conditions", model.AccessConditions)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Integrations", model.Integrations)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Credential Providers", model.CredentialProviders)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Server Workloads", model.ServerWorkloads)

	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Agent Controllers", model.AgentControllers)

	dto.Permissions = appendReadOnlyPermissionToDTO(dto.Permissions, "Access Authorization Events", model.AccessAuthorizationEvents)
	dto.Permissions = appendReadOnlyPermissionToDTO(dto.Permissions, "Audit Logs", model.AuditLogs)
	dto.Permissions = appendReadOnlyPermissionToDTO(dto.Permissions, "Workload Events", model.WorkloadEvents)

	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Users", model.Users)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Roles", model.Roles)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Log Streams", model.LogStreams)
	dto.Permissions = appendPermissionToDTO(dto.Permissions, "Identity Providers", model.IdentityProviders)

	return dto
}

func appendPermissionToDTO(list []aembit.RolePermissionDTO, name string, permission *rolePermission) []aembit.RolePermissionDTO {
	if permission == nil {
		return list
	}
	return append(list, aembit.RolePermissionDTO{
		Name:  name,
		Read:  permission.Read.ValueBool(),
		Write: permission.Write.ValueBool(),
	})
}

func appendReadOnlyPermissionToDTO(list []aembit.RolePermissionDTO, name string, permission *roleReadOnlyPermission) []aembit.RolePermissionDTO {
	if permission == nil {
		return list
	}
	return append(list, aembit.RolePermissionDTO{
		Name: name,
		Read: permission.Read.ValueBool(),
	})
}

// DTO to Model conversion methods.
func convertRoleDTOToModel(ctx context.Context, dto aembit.RoleDTO) roleResourceModel {
	var model roleResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	for _, permission := range dto.Permissions {
		switch permission.Name {
		case "Access Policies":
			model.AccessPolicies = convertPermissionDTOToPermission(permission)
		case "Client Workloads":
			model.ClientWorkloads = convertPermissionDTOToPermission(permission)
		case "Trust Providers":
			model.TrustProviders = convertPermissionDTOToPermission(permission)
		case "Access Conditions":
			model.AccessConditions = convertPermissionDTOToPermission(permission)
		case "Integrations":
			model.Integrations = convertPermissionDTOToPermission(permission)
		case "Credential Providers":
			model.CredentialProviders = convertPermissionDTOToPermission(permission)
		case "Server Workloads":
			model.ServerWorkloads = convertPermissionDTOToPermission(permission)
		case "Agent Controllers":
			model.AgentControllers = convertPermissionDTOToPermission(permission)
		case "Access Authorization Events":
			model.AccessAuthorizationEvents = convertPermissionDTOToReadOnlyPermission(permission)
		case "Audit Logs":
			model.AuditLogs = convertPermissionDTOToReadOnlyPermission(permission)
		case "Workload Events":
			model.WorkloadEvents = convertPermissionDTOToReadOnlyPermission(permission)
		case "Users":
			model.Users = convertPermissionDTOToPermission(permission)
		case "Roles":
			model.Roles = convertPermissionDTOToPermission(permission)
		case "Log Streams":
			model.LogStreams = convertPermissionDTOToPermission(permission)
		case "Identity Providers":
			model.IdentityProviders = convertPermissionDTOToPermission(permission)
		}
	}

	return model
}

func convertPermissionDTOToPermission(permission aembit.RolePermissionDTO) *rolePermission {
	return &rolePermission{
		Read:  types.BoolValue(permission.Read),
		Write: types.BoolValue(permission.Write),
	}
}

func convertPermissionDTOToReadOnlyPermission(permission aembit.RolePermissionDTO) *roleReadOnlyPermission {
	return &roleReadOnlyPermission{
		Read: types.BoolValue(permission.Read),
	}
}
