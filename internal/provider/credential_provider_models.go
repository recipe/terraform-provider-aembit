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
	Tags                   types.Map                                      `tfsdk:"tags"`
	AembitToken            *credentialProviderAembitTokenModel            `tfsdk:"aembit_access_token"`
	APIKey                 *credentialProviderAPIKeyModel                 `tfsdk:"api_key"`
	AwsSTS                 *credentialProviderAwsSTSModel                 `tfsdk:"aws_sts"`
	GoogleWorkload         *credentialProviderGoogleWorkloadModel         `tfsdk:"google_workload_identity"`
	SnowflakeToken         *credentialProviderSnowflakeTokenModel         `tfsdk:"snowflake_jwt"`
	OAuthClientCredentials *credentialProviderOAuthClientCredentialsModel `tfsdk:"oauth_client_credentials"`
	UsernamePassword       *credentialProviderUserPassModel               `tfsdk:"username_password"`
	VaultClientToken       *credentialProviderVaultClientTokenModel       `tfsdk:"vault_client_token"`
}

// credentialProviderDataSourceModel maps the datasource schema.
type credentialProvidersDataSourceModel struct {
	CredentialProviders []credentialProviderResourceModel `tfsdk:"credential_providers"`
}

type credentialProviderAembitTokenModel struct {
	Audience types.String `tfsdk:"audience"`
	Role     types.String `tfsdk:"role_id"`
	Lifetime int32        `tfsdk:"lifetime"`
}

type credentialProviderAPIKeyModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

type credentialProviderAwsSTSModel struct {
	OIDCIssuer    types.String `tfsdk:"oidc_issuer"`
	RoleARN       types.String `tfsdk:"role_arn"`
	TokenAudience types.String `tfsdk:"token_audience"`
	Lifetime      int32        `tfsdk:"lifetime"`
}

type credentialProviderGoogleWorkloadModel struct {
	OIDCIssuer     types.String `tfsdk:"oidc_issuer"`
	Audience       types.String `tfsdk:"audience"`
	ServiceAccount types.String `tfsdk:"service_account"`
	Lifetime       int32        `tfsdk:"lifetime"`
}

type credentialProviderSnowflakeTokenModel struct {
	AccountID        types.String `tfsdk:"account_id"`
	Username         types.String `tfsdk:"username"`
	AlertUserCommand types.String `tfsdk:"alter_user_command"`
}

// credentialProviderOAuthClientCredentialsModel maps OAuth Client Credentials Flow configuration.
type credentialProviderOAuthClientCredentialsModel struct {
	TokenURL     types.String `tfsdk:"token_url"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Scopes       types.String `tfsdk:"scopes"`
}

type credentialProviderUserPassModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// credentialProviderVaultClientTokenModel maps OAuth Client Credentials Flow configuration.
type credentialProviderVaultClientTokenModel struct {
	Subject         string                                                 `tfsdk:"subject"`
	SubjectType     string                                                 `tfsdk:"subject_type"`
	CustomClaims    []*credentialProviderVaultClientTokenCustomClaimsModel `tfsdk:"custom_claims"`
	Lifetime        int32                                                  `tfsdk:"lifetime"`
	VaultHost       string                                                 `tfsdk:"vault_host"`
	VaultTLS        bool                                                   `tfsdk:"vault_tls"`
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
