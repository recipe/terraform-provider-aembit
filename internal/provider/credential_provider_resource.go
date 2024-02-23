package provider

import (
	"context"
	"encoding/json"
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
	_ resource.Resource                = &credentialProviderResource{}
	_ resource.ResourceWithConfigure   = &credentialProviderResource{}
	_ resource.ResourceWithImportState = &credentialProviderResource{}
)

// NewCredentialProviderResource is a helper function to simplify the provider implementation.
func NewCredentialProviderResource() resource.Resource {
	return &credentialProviderResource{}
}

// credentialProviderResource is the resource implementation.
type credentialProviderResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *credentialProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_provider"
}

// Configure adds the provider configured client to the resource.
func (r *credentialProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.CloudClient)

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
func (r *credentialProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the credential provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the credential provider.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the credential provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the credential provider.",
				Optional:    true,
				Computed:    true,
			},
			"api_key": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"api_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
				},
			},
			"oauth_client_credentials": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"token_url": schema.StringAttribute{
						Required: true,
					},
					"client_id": schema.StringAttribute{
						Required: true,
					},
					"client_secret": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"scopes": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"vault_client_token": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"subject": schema.StringAttribute{
						Required: true,
					},
					"subject_type": schema.StringAttribute{
						Required: true,
					},
					"custom_claims": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Required: true,
								},
								"value": schema.StringAttribute{
									Required: true,
								},
								"value_type": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"lifetime": schema.Int64Attribute{
						Required: true,
					},
					"vault_host": schema.StringAttribute{
						Required: true,
					},
					"vault_port": schema.Int64Attribute{
						Required: true,
					},
					"vault_tls": schema.BoolAttribute{
						Required: true,
					},
					"vault_namespace": schema.StringAttribute{
						Optional: true,
					},
					"vault_role": schema.StringAttribute{
						Optional: true,
					},
					"vault_path": schema.StringAttribute{
						Required: true,
					},
					"vault_forwarding": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one credential provider type is specified.
func (r *credentialProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("api_key"),
			path.MatchRoot("oauth_client_credentials"),
			path.MatchRoot("vault_client_token"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *credentialProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan credentialProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var credential aembit.CredentialProviderDTO = convertCredentialProviderModelToDTO(ctx, plan, nil, r.client.Tenant, r.client.StackDomain)

	// Create new Credential Provider
	credentialProvider, err := r.client.CreateCredentialProvider(credential, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating credential provider",
			"Could not create credential provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderDTOToModel(*credentialProvider, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *credentialProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed credential value from Aembit
	credentialProvider, err := r.client.GetCredentialProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Credential Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertCredentialProviderDTOToModel(credentialProvider, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *credentialProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	var externalID string = state.ID.ValueString()

	// Retrieve values from plan
	var plan credentialProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var credential aembit.CredentialProviderDTO = convertCredentialProviderModelToDTO(ctx, plan, &externalID, r.client.Tenant, r.client.StackDomain)

	// Update Credential Provider
	credentialProvider, err := r.client.UpdateCredentialProvider(credential, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating credential provider",
			"Could not update credential provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderDTOToModel(*credentialProvider, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *credentialProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Credential Provider is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Credential Provider is active and cannot be deleted. Please mark the credential as inactive first.",
		)
		return
	}

	// Delete existing Credential Provider
	_, err := r.client.DeleteCredentialProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Could not delete credential provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *credentialProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertCredentialProviderModelToDTO(_ context.Context, model credentialProviderResourceModel, externalID *string, tenantID string, stackDomain string) aembit.CredentialProviderDTO {
	var credential aembit.CredentialProviderDTO
	credential.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if externalID != nil {
		credential.EntityDTO.ExternalID = *externalID
	}

	// Handle the API Key use case
	if model.APIKey != nil {
		credential.Type = "apikey"
		apiKey := aembit.CredentialAPIKeyDTO{APIKey: model.APIKey.APIKey.ValueString()}
		apiKeyJSON, _ := json.Marshal(apiKey)
		credential.ProviderDetail = string(apiKeyJSON)
	}

	// Handle the OAuth Client Credentials use case
	if model.OAuthClientCredentials != nil {
		credential.Type = "oauth-client-credential"
		oauth := aembit.CredentialOAuthClientCredentialDTO{
			TokenURL:        model.OAuthClientCredentials.TokenURL.ValueString(),
			ClientID:        model.OAuthClientCredentials.ClientID.ValueString(),
			ClientSecret:    model.OAuthClientCredentials.ClientSecret.ValueString(),
			Scope:           model.OAuthClientCredentials.Scopes.ValueString(),
			CredentialStyle: "authHeader",
		}
		oauthJSON, _ := json.Marshal(oauth)
		credential.ProviderDetail = string(oauthJSON)
	}

	// Handle the Vault Client Token use case
	if model.VaultClientToken != nil {
		credential.Type = "vaultClientToken"

		vault := aembit.CredentialVaultClientTokenDTO{
			JwtConfig: &aembit.CredentialVaultClientTokenJwtConfigDTO{
				Issuer:       fmt.Sprintf("https://%s.id.%s/", tenantID, stackDomain),
				Subject:      model.VaultClientToken.Subject,
				SubjectType:  model.VaultClientToken.SubjectType,
				Lifetime:     model.VaultClientToken.Lifetime,
				CustomClaims: make([]aembit.CredentialVaultClientTokenClaimsDTO, len(model.VaultClientToken.CustomClaims)),
			},
			VaultCluster: &aembit.CredentialVaultClientTokenVaultClusterDTO{
				VaultHost:          model.VaultClientToken.VaultHost,
				Port:               model.VaultClientToken.VaultPort,
				TLS:                model.VaultClientToken.VaultTLS,
				Namespace:          model.VaultClientToken.VaultNamespace,
				Role:               model.VaultClientToken.VaultRole,
				AuthenticationPath: model.VaultClientToken.VaultPath,
				ForwardingConfig:   model.VaultClientToken.VaultForwarding,
			},
		}
		for i, claim := range model.VaultClientToken.CustomClaims {
			vault.JwtConfig.CustomClaims[i] = aembit.CredentialVaultClientTokenClaimsDTO{
				Key:       claim.Key,
				Value:     claim.Value,
				ValueType: claim.ValueType,
			}
		}

		vaultJSON, _ := json.Marshal(vault)
		credential.ProviderDetail = string(vaultJSON)
	}
	return credential
}

func convertCredentialProviderDTOToModel(dto aembit.CredentialProviderDTO, state credentialProviderResourceModel) credentialProviderResourceModel {
	var model credentialProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	// Set the objects to null to begin with
	model.APIKey = nil
	model.OAuthClientCredentials = nil
	model.VaultClientToken = nil

	// Now fill in the objects based on the Credential Provider type
	switch dto.Type {
	case "apikey":
		model.APIKey = convertAPIKeyDTOToModel(dto, state)
	case "oauth-client-credential":
		model.OAuthClientCredentials = convertOAuthClientCredentialDTOToModel(dto, state)
	case "vaultClientToken":
		model.VaultClientToken = convertVaultClientTokenDTOToModel(dto, state)
	}
	return model
}

// convertAPIKeyDTOToModel converts the API Key state object into a model ready for terraform processing.
// Note: Since Aembit vaults the API Key and does not return it in the API, the DTO will never contain the stored value.
func convertAPIKeyDTOToModel(_ aembit.CredentialProviderDTO, state credentialProviderResourceModel) *credentialProviderAPIKeyModel {
	value := credentialProviderAPIKeyModel{APIKey: types.StringNull()}
	if state.APIKey != nil {
		value.APIKey = state.APIKey.APIKey
	}
	return &value
}

// convertOAuthClientCredentialDTOToModel converts the OAuth Client Credential state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Client Secret and does not return it in the API, the DTO will never contain the stored value.
func convertOAuthClientCredentialDTOToModel(dto aembit.CredentialProviderDTO, state credentialProviderResourceModel) *credentialProviderOAuthClientCredentialsModel {
	value := credentialProviderOAuthClientCredentialsModel{ClientSecret: types.StringNull()}

	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var oauth aembit.CredentialOAuthClientCredentialDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &oauth)
	if err != nil {
		return nil
	}

	value.TokenURL = types.StringValue(oauth.TokenURL)
	value.ClientID = types.StringValue(oauth.ClientID)
	value.Scopes = types.StringValue(oauth.Scope)
	if state.OAuthClientCredentials != nil {
		value.ClientSecret = state.OAuthClientCredentials.ClientSecret
	}

	return &value
}

// convertVaultClientTokenDTOToModel converts the VaultClientToken state object into a model ready for terraform processing.
func convertVaultClientTokenDTOToModel(dto aembit.CredentialProviderDTO, _ credentialProviderResourceModel) *credentialProviderVaultClientTokenModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var vault aembit.CredentialVaultClientTokenDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &vault)
	if err != nil {
		return nil
	}

	value := credentialProviderVaultClientTokenModel{
		Subject:     vault.JwtConfig.Subject,
		SubjectType: vault.JwtConfig.SubjectType,
		Lifetime:    vault.JwtConfig.Lifetime,

		VaultHost:       vault.VaultCluster.VaultHost,
		VaultPort:       vault.VaultCluster.Port,
		VaultTLS:        vault.VaultCluster.TLS,
		VaultNamespace:  vault.VaultCluster.Namespace,
		VaultRole:       vault.VaultCluster.Role,
		VaultPath:       vault.VaultCluster.AuthenticationPath,
		VaultForwarding: vault.VaultCluster.ForwardingConfig,
	}

	// Get the custom claims to be injected into the model
	claims := make([]*credentialProviderVaultClientTokenCustomClaimsModel, len(vault.JwtConfig.CustomClaims))
	//types.ObjectValue(credentialProviderVaultClientTokenCustomClaimsModel.AttrTypes),
	//claims := getSetObjectAttr(ctx, model.VaultClientToken, "custom_claims")
	for i, claim := range vault.JwtConfig.CustomClaims {
		claims[i] = &credentialProviderVaultClientTokenCustomClaimsModel{
			Key:       claim.Key,
			Value:     claim.Value,
			ValueType: claim.ValueType,
		}
	}
	value.CustomClaims = claims
	return &value
}
