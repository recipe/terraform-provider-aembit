// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure AembitProvider satisfies various provider interfaces.
var _ provider.Provider = &aembitProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &aembitProvider{
			version: version,
		}
	}
}

// aembitProviderModel maps provider schema data to a Go type.
type aembitProviderModel struct {
	Tenant      types.String `tfsdk:"tenant"`
	Token       types.String `tfsdk:"token"`
	StackDomain types.String `tfsdk:"stack_domain"`
}

// AembitProvider defines the provider implementation.
type aembitProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *aembitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aembit"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *aembitProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"stack_domain": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *aembitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Aembit client")

	// Retrieve provider data from configuration
	var config aembitProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Tenant.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant"),
			"Unknown Aembit API Tenant",
			"The provider cannot create the Aembit API client as there is an unknown configuration value for the Aembit API Tenant. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AEMBIT_TENANT_ID environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Aembit API Access Token",
			"The provider cannot create the Aembit API client as there is an unknown configuration value for the Aembit API Access Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AEMBIT_TOKEN environment variable.",
		)
	}

	if config.StackDomain.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("stack_domain"),
			"Unknown Aembit API Stack Domain",
			"The provider cannot create the Aembit API client as there is an unknown configuration value for the Aembit API Stack Domain. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AEMBIT_STACK_DOMAIN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	tenant := os.Getenv("AEMBIT_TENANT_ID")
	token := os.Getenv("AEMBIT_TOKEN")
	stackDomain := os.Getenv("AEMBIT_STACK_DOMAIN")

	if !config.Tenant.IsNull() {
		tenant = config.Tenant.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.StackDomain.IsNull() {
		stackDomain = config.StackDomain.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if tenant == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant"),
			"Missing Aembit API Tenant",
			"The provider cannot create the Aembit API client as there is a missing or empty value for the Aembit API Tenant. "+
				"Set the host value in the configuration or use the AEMBIT_TENANT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Aembit API Access Token",
			"The provider cannot create the Aembit API client as there is a missing or empty value for the Aembit API Access Token. "+
				"Set the password value in the configuration or use the AEMBIT_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if stackDomain == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("stack_domain"),
			"Missing Aembit API Stack Domain",
			"The provider cannot create the Aembit API client as there is a missing or empty value for the Aembit API Stack Domain. "+
				"Set the password value in the configuration or use the AEMBIT_STACK_DOMAIN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "aembit_tenant", tenant)
	ctx = tflog.SetField(ctx, "aembit_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "aembit_token")

	tflog.Debug(ctx, "Creating Aembit client")

	// Create a new Aembit client using the configuration values
	client, err := aembit.NewClient(aembit.URLBuilder{}, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Aembit API Client",
			"An unexpected error occurred when creating the Aembit API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Aembit Client Error: "+err.Error(),
		)
		return
	}
	client.Tenant = tenant
	client.StackDomain = stackDomain

	// Make the Aembit client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Aembit client", map[string]any{"success": true})
}

func (p *aembitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerWorkloadResource,
		NewCredentialProviderResource,
		NewTrustProviderResource,
		NewClientWorkloadResource,
		NewIntegrationResource,
		NewAccessConditionResource,
		NewAccessPolicyResource,
	}
}

func (p *aembitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServerWorkloadsDataSource,
		NewCredentialProvidersDataSource,
		NewTrustProvidersDataSource,
		NewClientWorkloadsDataSource,
		NewAccessPoliciesDataSource,
	}
}
