package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// credentialProviderResourceModel maps the resource schema.
type credentialProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                     types.String                                   `tfsdk:"id"`
	Name                   types.String                                   `tfsdk:"name"`
	Description            types.String                                   `tfsdk:"description"`
	IsActive               types.Bool                                     `tfsdk:"is_active"`
	ApiKey                 *credentialProviderApiKeyModel                 `tfsdk:"api_key"`
	OAuthClientCredentials *credentialProviderOAuthClientCredentialsModel `tfsdk:"oauth_client_credentials"`
	VaultClientToken       *credentialProviderVaultClientTokenModel       `tfsdk:"vault_client_token"`
}

// credentialProviderDataSourceModel maps the datasource schema.
type credentialProvidersDataSourceModel struct {
	CredentialProviders []credentialProviderResourceModel `tfsdk:"credential_providers"`
}

type credentialProviderApiKeyModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

// credentialProviderOAuthClientCredentialsModel maps OAuth Client Credentials Flow configuration.
type credentialProviderOAuthClientCredentialsModel struct {
	TokenUrl     types.String `tfsdk:"token_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Scopes       types.String `tfsdk:"scopes"`
}

// credentialProviderVaultClientTokenModel maps OAuth Client Credentials Flow configuration.
type credentialProviderVaultClientTokenModel struct {
	Subject         string                                                 `tfsdk:"subject"`
	SubjectType     string                                                 `tfsdk:"subject_type"`
	CustomClaims    []*credentialProviderVaultClientTokenCustomClaimsModel `tfsdk:"custom_claims"`
	Lifetime        int32                                                  `tfsdk:"lifetime"`
	VaultHost       string                                                 `tfsdk:"vault_host"`
	VaultTls        bool                                                   `tfsdk:"vault_tls"`
	VaultPort       int32                                                  `tfsdk:"vault_port"`
	VaultNamespace  string                                                 `tfsdk:"vault_namespace"`
	VaultRole       string                                                 `tfsdk:"vault_role"`
	VaultPath       string                                                 `tfsdk:"vault_path"`
	VaultForwarding string                                                 `tfsdk:"vault_forwarding"`
}

type credentialProviderVaultClientTokenCustomClaimsModel struct {
	Key       string `tfsdk:"key"`
	Value     string `tfsdk:"value"`
	ValueType string `tfsdk:"value_type"`
}
