package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &credentialProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialProvidersDataSource{}
)

// NewCredentialProvidersDataSource is a helper function to simplify the provider implementation.
func NewCredentialProvidersDataSource() datasource.DataSource {
	return &credentialProvidersDataSource{}
}

// credentialProvidersDataSource is the data source implementation.
type credentialProvidersDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *credentialProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Metadata returns the data source type name.
func (d *credentialProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_providers"
}

// Schema defines the schema for the resource.
func (d *credentialProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an credential provider.",
		Attributes: map[string]schema.Attribute{
			"credential_providers": schema.ListNestedAttribute{
				Description: "List of credential providers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the credential provider.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the credential provider.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the credential provider.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the credential provider.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Tags are key-value pairs.",
							ElementType: types.StringType,
							Computed:    true,
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
									Computed:    true,
								},
								"lifetime": schema.Int64Attribute{
									Description: "Lifetime of the Credential Provider.",
									Computed:    true,
								},
							},
						},
						"api_key": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"api_key": schema.StringAttribute{
									Computed:  true,
									Sensitive: true,
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
									Computed:    true,
								},
								"token_audience": schema.StringAttribute{
									Description: "Token Audience for AWS IAM Identity Provider configuration of the Credential Provider.",
									Computed:    true,
								},
								"lifetime": schema.Int64Attribute{
									Description: "Lifetime (seconds) of the AWS Session credentials requested by the Credential Provider.",
									Optional:    true,
									Computed:    true,
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
									Computed:    true,
								},
								"service_account": schema.StringAttribute{
									Description: "Service Account email of the GCP Session credentials requested by the Credential Provider.",
									Computed:    true,
								},
								"lifetime": schema.Int64Attribute{
									Description: "Lifetime (seconds) of the GCP Session credentials requested by the Credential Provider.",
									Optional:    true,
									Computed:    true,
								},
							},
						},
						"snowflake_jwt": schema.SingleNestedAttribute{
							Description: "JSON Web Token type Credential Provider configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Description: "Snowflake Account ID of the Credential Provider.",
									Computed:    true,
								},
								"username": schema.StringAttribute{
									Description: "Snowflake Username of the Credential Provider.",
									Computed:    true,
								},
								"alter_user_command": schema.StringAttribute{
									Description: "Snowflake Alter User Command generated for configuration of Snowflake by the Credential Provider.",
									Computed:    true,
								},
							},
						},
						"oauth_client_credentials": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"token_url": schema.StringAttribute{
									Computed: true,
								},
								"client_id": schema.StringAttribute{
									Computed: true,
								},
								"client_secret": schema.StringAttribute{
									Computed:  true,
									Sensitive: true,
								},
								"scopes": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"username_password": schema.SingleNestedAttribute{
							Description: "Username/Password type Credential Provider configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"username": schema.StringAttribute{
									Description: "Username of the Credential Provider.",
									Computed:    true,
								},
								"password": schema.StringAttribute{
									Description: "Password of the Credential Provider.",
									Computed:    true,
									Sensitive:   true,
								},
							},
						},
						"vault_client_token": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"subject": schema.StringAttribute{
									Computed: true,
								},
								"subject_type": schema.StringAttribute{
									Computed: true,
								},
								"custom_claims": schema.SetNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Computed: true,
											},
											"value": schema.StringAttribute{
												Computed: true,
											},
											"value_type": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
								"lifetime": schema.Int64Attribute{
									Computed: true,
								},
								"vault_host": schema.StringAttribute{
									Computed: true,
								},
								"vault_port": schema.Int64Attribute{
									Computed: true,
								},
								"vault_tls": schema.BoolAttribute{
									Computed: true,
								},
								"vault_namespace": schema.StringAttribute{
									Computed: true,
								},
								"vault_role": schema.StringAttribute{
									Computed: true,
								},
								"vault_path": schema.StringAttribute{
									Computed: true,
								},
								"vault_forwarding": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *credentialProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state credentialProvidersDataSourceModel

	credentialProviders, err := d.client.GetCredentialProviders(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Credential Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, credentialProvider := range credentialProviders {
		credentialProviderState := convertCredentialProviderDTOToModel(ctx, credentialProvider, credentialProviderResourceModel{}, d.client.Tenant, d.client.StackDomain)
		state.CredentialProviders = append(state.CredentialProviders, credentialProviderState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
