package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// integrationResourceModel maps the resource schema.
type integrationResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                     types.String                            `tfsdk:"id"`
	Name                   types.String                            `tfsdk:"name"`
	Description            types.String                            `tfsdk:"description"`
	IsActive               types.Bool                              `tfsdk:"is_active"`
	Type                   types.String                            `tfsdk:"type"`
	SyncFrequency          types.Int64                             `tfsdk:"sync_frequency"`
	Endpoint               types.String                            `tfsdk:"endpoint"`
	OAuthClientCredentials *integrationOAuthClientCredentialsModel `tfsdk:"oauth_client_credentials"`
}

type integrationOAuthClientCredentialsModel struct {
	TokenURL     types.String `tfsdk:"token_url"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Audience     types.String `tfsdk:"audience"`
}

// integrationDataSourceModel maps the datasource schema.
type integrationsDataSourceModel struct {
	Integrations []integrationResourceModel `tfsdk:"integrations"`
}
