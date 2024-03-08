package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
			fmt.Sprintf("Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
				Description: "Unique identifier of the Credential Provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Credential Provider.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Credential Provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Credential Provider.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"aembit_access_token": schema.SingleNestedAttribute{
				Description: "Aembit Access Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"audience": schema.StringAttribute{
						Description: "Audience of the Credential Provider.",
						Computed:    true,
					},
					"role_id": schema.StringAttribute{
						Description: "Aembit Role ID of the Credential Provider.",
						Required:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime of the Credential Provider.",
						Required:    true,
					},
				},
			},
			"api_key": schema.SingleNestedAttribute{
				Description: "API Key type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"api_key": schema.StringAttribute{
						Description: "API Key secret of the Credential Provider.",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
			"aws_sts": schema.SingleNestedAttribute{
				Description: "AWS Security Token Service Federation type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_issuer": schema.StringAttribute{
						Description: "OIDC Issuer for AWS IAM Identity Provider configuration of the Credential Provider.",
						Computed:    true,
					},
					"role_arn": schema.StringAttribute{
						Description: "AWS Role Arn to be used for the AWS Session credentials requested by the Credential Provider.",
						Required:    true,
					},
					"token_audience": schema.StringAttribute{
						Description: "Token Audience for AWS IAM Identity Provider configuration of the Credential Provider.",
						Computed:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (seconds) of the AWS Session credentials requested by the Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3600),
					},
				},
			},
			"google_workload_identity": schema.SingleNestedAttribute{
				Description: "Google Workload Identity Federation type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_issuer": schema.StringAttribute{
						Description: "OIDC Issuer for AWS IAM Identity Provider configuration of the Credential Provider.",
						Computed:    true,
					},
					"audience": schema.StringAttribute{
						Description: "Audience for GCP Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
					},
					"service_account": schema.StringAttribute{
						Description: "Service Account email of the GCP Session credentials requested by the Credential Provider.",
						Required:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (seconds) of the GCP Session credentials requested by the Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3600),
					},
				},
			},
			"snowflake_jwt": schema.SingleNestedAttribute{
				Description: "JSON Web Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "Snowflake Account ID of the Credential Provider.",
						Required:    true,
					},
					"username": schema.StringAttribute{
						Description: "Snowflake Username of the Credential Provider.",
						Required:    true,
					},
					"alter_user_command": schema.StringAttribute{
						Description: "Snowflake Alter User Command generated for configuration of Snowflake by the Credential Provider.",
						Computed:    true,
					},
				},
			},
			"oauth_client_credentials": schema.SingleNestedAttribute{
				Description: "OAuth Client Credentials Flow type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"token_url": schema.StringAttribute{
						Description: "Token URL for the OAuth Credential Provider.",
						Required:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "Client ID for the OAuth Credential Provider.",
						Required:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Client Secret for the OAuth Credential Provider.",
						Optional:    true,
						Sensitive:   true,
					},
					"scopes": schema.StringAttribute{
						Description: "Scopes for the OAuth Credential Provider.",
						Optional:    true,
					},
				},
			},
			"username_password": schema.SingleNestedAttribute{
				Description: "Username/Password type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Description: "Username of the Credential Provider.",
						Optional:    true,
					},
					"password": schema.StringAttribute{
						Description: "Password of the Credential Provider.",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
			"vault_client_token": schema.SingleNestedAttribute{
				Description: "Vault Client Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"subject": schema.StringAttribute{
						Description: "Subject of the JWT Token used to authenticate to the Vault Cluster.",
						Required:    true,
					},
					"subject_type": schema.StringAttribute{
						Description: "Type of value for the JWT Token Subject. Possible values are `literal` or `dynamic`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"literal",
								"dynamic",
							}...),
						},
					},
					"custom_claims": schema.SetNestedAttribute{
						Description: "Set of Custom Claims for the JWT Token used to authenticate to the Vault Cluster.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: "Key for the JWT Token Custom Claim.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "Value for the JWT Token Custom Claim.",
									Required:    true,
								},
								"value_type": schema.StringAttribute{
									Description: "Type of value for the JWT Token Custom Claim. Possible values are `literal` or `dynamic`.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{
											"literal",
											"dynamic",
										}...),
									},
								},
							},
						},
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime of the JWT Token used to authenticate to the Vault Cluster. Note: The lifetime of the retrieved Vault Client Token is managed within Vault configuration.",
						Required:    true,
					},
					"vault_host": schema.StringAttribute{
						Description: "Hostname of the Vault Cluster to be used for executing the login API.",
						Required:    true,
					},
					"vault_port": schema.Int64Attribute{
						Description: "Port of the Vault Cluster to be used for executing the login API.",
						Required:    true,
					},
					"vault_tls": schema.BoolAttribute{
						Description: "Configuration to utilize TLS for connectivity to the Vault Cluster.",
						Required:    true,
					},
					"vault_namespace": schema.StringAttribute{
						Description: "Namespace to utilize when executing the login API on the Vault Cluster.",
						Optional:    true,
					},
					"vault_role": schema.StringAttribute{
						Description: "Role to utilize when executing the login API on the Vault Cluster.",
						Optional:    true,
					},
					"vault_path": schema.StringAttribute{
						Description: "Path to utilize when executing the login API on the Vault Cluster.",
						Required:    true,
					},
					"vault_forwarding": schema.StringAttribute{
						Description: "If Vault Forwarding is required, this configuration can be set to `unconditional` or `conditional`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"",
								"unconditional",
								"conditional",
							}...),
						},
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one Credential Provider type is specified.
func (r *credentialProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("aembit_access_token"),
			path.MatchRoot("api_key"),
			path.MatchRoot("aws_sts"),
			path.MatchRoot("google_workload_identity"),
			path.MatchRoot("snowflake_jwt"),
			path.MatchRoot("oauth_client_credentials"),
			path.MatchRoot("username_password"),
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
			"Error creating Credential Provider",
			"Could not create Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderDTOToModel(ctx, *credentialProvider, plan, r.client.Tenant, r.client.StackDomain)

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

	state = convertCredentialProviderDTOToModel(ctx, credentialProvider, state, r.client.Tenant, r.client.StackDomain)

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
			"Error updating Credential Provider",
			"Could not update Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderDTOToModel(ctx, *credentialProvider, plan, r.client.Tenant, r.client.StackDomain)

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

	// Check if Credential Provider is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableCredentialProvider(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Credential Provider",
				"Could not disable Credential Provider, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Credential Provider
	_, err := r.client.DeleteCredentialProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Could not delete Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *credentialProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertCredentialProviderModelToDTO(ctx context.Context, model credentialProviderResourceModel, externalID *string, tenantID string, stackDomain string) aembit.CredentialProviderDTO {
	var credential aembit.CredentialProviderDTO
	credential.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			credential.Tags = append(credential.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}
	if externalID != nil {
		credential.EntityDTO.ExternalID = *externalID
	}

	// Handle the Aembit Token use case
	if model.AembitToken != nil {
		credential.Type = "aembit-access-token"
		aembitToken := aembit.CredentialAembitTokenDTO{
			Audience: fmt.Sprintf("%s.api.%s", tenantID, stackDomain),
			RoleID:   model.AembitToken.Role.ValueString(),
			Lifetime: model.AembitToken.Lifetime,
		}
		aembitTokenJSON, _ := json.Marshal(aembitToken)
		credential.ProviderDetail = string(aembitTokenJSON)
	}

	// Handle the API Key use case
	if model.APIKey != nil {
		credential.Type = "apikey"
		apiKey := aembit.CredentialAPIKeyDTO{APIKey: model.APIKey.APIKey.ValueString()}
		apiKeyJSON, _ := json.Marshal(apiKey)
		credential.ProviderDetail = string(apiKeyJSON)
	}

	// Handle the AWS STS use case
	if model.AwsSTS != nil {
		credential.Type = "aws-sts-oidc"
		awsSTS := aembit.CredentialAwsSTSDTO{
			RoleArn:  model.AwsSTS.RoleARN.ValueString(),
			Lifetime: model.AwsSTS.Lifetime,
		}
		awsSTSJSON, _ := json.Marshal(awsSTS)
		credential.ProviderDetail = string(awsSTSJSON)
	}

	// Handle the GCP Workload Identity Federation use case
	if model.GoogleWorkload != nil {
		credential.Type = "gcp-identity-federation"
		gcpWorkload := aembit.CredentialGoogleWorkloadDTO{
			Audience:       model.GoogleWorkload.Audience.ValueString(),
			ServiceAccount: model.GoogleWorkload.ServiceAccount.ValueString(),
			Lifetime:       model.GoogleWorkload.Lifetime,
		}
		gcpWorkloadJSON, _ := json.Marshal(gcpWorkload)
		credential.ProviderDetail = string(gcpWorkloadJSON)
	}

	// Handle the Snowflake JWT use case
	if model.SnowflakeToken != nil {
		credential.Type = "signed-jwt"
		gcpWorkload := aembit.CredentialSnowflakeTokenDTO{
			TokenConfiguration: "snowflake",
			AlgorithmType:      "RS256",
			Issuer:             fmt.Sprintf("%s.%s.SHA256:{sha256(publicKey)}", model.SnowflakeToken.AccountID.ValueString(), model.SnowflakeToken.Username.ValueString()),
			Subject:            fmt.Sprintf("%s.%s", model.SnowflakeToken.AccountID.ValueString(), model.SnowflakeToken.Username.ValueString()),
			Lifetime:           1,
		}
		gcpWorkloadJSON, _ := json.Marshal(gcpWorkload)
		credential.ProviderDetail = string(gcpWorkloadJSON)
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

	// Handle the Username Password use case
	if model.UsernamePassword != nil {
		credential.Type = "username-password"
		userPass := aembit.CredentialUsernamePasswordDTO{
			Username: model.UsernamePassword.Username.ValueString(),
			Password: model.UsernamePassword.Password.ValueString(),
		}
		userPassJSON, _ := json.Marshal(userPass)
		credential.ProviderDetail = string(userPassJSON)
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

func convertCredentialProviderDTOToModel(ctx context.Context, dto aembit.CredentialProviderDTO, state credentialProviderResourceModel, tenant, stackDomain string) credentialProviderResourceModel {
	var model credentialProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	// Set the objects to null to begin with
	model.AembitToken = nil
	model.APIKey = nil
	model.AwsSTS = nil
	model.GoogleWorkload = nil
	model.OAuthClientCredentials = nil
	model.UsernamePassword = nil
	model.VaultClientToken = nil

	// Now fill in the objects based on the Credential Provider type
	switch dto.Type {
	case "aembit-access-token":
		model.AembitToken = convertAembitTokenDTOToModel(dto)
	case "apikey":
		model.APIKey = convertAPIKeyDTOToModel(dto, state)
	case "aws-sts-oidc":
		model.AwsSTS = convertAwsSTSDTOToModel(dto, tenant, stackDomain)
	case "gcp-identity-federation":
		model.GoogleWorkload = convertGoogleWorkloadDTOToModel(dto, tenant, stackDomain)
	case "signed-jwt":
		model.SnowflakeToken = convertSnowflakeTokenDTOToModel(dto)
	case "oauth-client-credential":
		model.OAuthClientCredentials = convertOAuthClientCredentialDTOToModel(dto, state)
	case "username-password":
		model.UsernamePassword = convertUserPassDTOToModel(dto, state)
	case "vaultClientToken":
		model.VaultClientToken = convertVaultClientTokenDTOToModel(dto, state)
	}
	return model
}

// convertAembitTokenDTOToModel converts the Aembit Token state object into a model ready for terraform processing.
func convertAembitTokenDTOToModel(dto aembit.CredentialProviderDTO) *credentialProviderAembitTokenModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var aembitToken aembit.CredentialAembitTokenDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &aembitToken)
	if err != nil {
		return nil
	}

	value := credentialProviderAembitTokenModel{
		Audience: types.StringValue(aembitToken.Audience),
		Role:     types.StringValue(aembitToken.RoleID),
		Lifetime: aembitToken.Lifetime,
	}
	return &value
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

// convertAwsSTSDTOToModel converts the AWS STS state object into a model ready for terraform processing.
func convertAwsSTSDTOToModel(dto aembit.CredentialProviderDTO, tenant, stackDomain string) *credentialProviderAwsSTSModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var awsSTS aembit.CredentialAwsSTSDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &awsSTS)
	if err != nil {
		return nil
	}

	value := credentialProviderAwsSTSModel{
		OIDCIssuer:    types.StringValue(fmt.Sprintf("https://%s.id.%s", tenant, stackDomain)),
		TokenAudience: types.StringValue("sts.amazonaws.com"),
		RoleARN:       types.StringValue(awsSTS.RoleArn),
		Lifetime:      awsSTS.Lifetime,
	}
	return &value
}

// convertGoogleWorkloadDTOToModel converts the Google Workload state object into a model ready for terraform processing.
func convertGoogleWorkloadDTOToModel(dto aembit.CredentialProviderDTO, tenant, stackDomain string) *credentialProviderGoogleWorkloadModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var gcpWorkload aembit.CredentialGoogleWorkloadDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &gcpWorkload)
	if err != nil {
		return nil
	}

	value := credentialProviderGoogleWorkloadModel{
		OIDCIssuer:     types.StringValue(fmt.Sprintf("https://%s.id.%s", tenant, stackDomain)),
		Audience:       types.StringValue(gcpWorkload.Audience),
		ServiceAccount: types.StringValue(gcpWorkload.ServiceAccount),
		Lifetime:       gcpWorkload.Lifetime,
	}
	return &value
}

// convertSnowflakeTokenDTOToModel converts the Snowflake JWT Token state object into a model ready for terraform processing.
func convertSnowflakeTokenDTOToModel(dto aembit.CredentialProviderDTO) *credentialProviderSnowflakeTokenModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var snowflake aembit.CredentialSnowflakeTokenDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &snowflake)
	if err != nil {
		return nil
	}

	acctData := strings.Split(snowflake.Subject, ".")
	keyData := strings.ReplaceAll(snowflake.KeyContent, "\n", "")
	keyData = strings.Replace(keyData, "-----BEGIN PUBLIC KEY-----", "", 1)
	keyData = strings.Replace(keyData, "-----END PUBLIC KEY-----", "", 1)
	value := credentialProviderSnowflakeTokenModel{
		AccountID:        types.StringValue(acctData[0]),
		Username:         types.StringValue(acctData[1]),
		AlertUserCommand: types.StringValue(fmt.Sprintf("ALTER USER %s SET RSA_PUBLIC_KEY='%s'", acctData[1], keyData)),
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

// convertUserPassDTOToModel converts the API Key state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Password and does not return it in the API, the DTO will never contain the stored value.
func convertUserPassDTOToModel(dto aembit.CredentialProviderDTO, state credentialProviderResourceModel) *credentialProviderUserPassModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	var userPass aembit.CredentialUsernamePasswordDTO
	err := json.Unmarshal([]byte(dto.ProviderDetail), &userPass)
	if err != nil {
		return nil
	}

	value := credentialProviderUserPassModel{
		Username: types.StringValue(userPass.Username),
		Password: types.StringNull(),
	}
	if state.UsernamePassword != nil {
		value.Password = state.UsernamePassword.Password
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
