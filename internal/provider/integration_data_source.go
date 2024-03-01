package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &integrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &integrationsDataSource{}
)

// NewIntegrationsDataSource is a helper function to simplify the provider implementation.
func NewIntegrationsDataSource() datasource.DataSource {
	return &integrationsDataSource{}
}

// integrationsDataSource is the data source implementation.
type integrationsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *integrationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *integrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integrations"
}

// Schema defines the schema for the resource.
func (d *integrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an integration.",
		Attributes: map[string]schema.Attribute{
			"integrations": schema.ListNestedAttribute{
				Description: "List of integrations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the integration.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the integration.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the integration.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the integration.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Tags are key-value pairs.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of Aembit integration (either `WizIntegrationApi` or `CrowdStrike`).",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"WizIntegrationApi", "CrowdStrike"}...),
							},
						},
						"sync_frequency": schema.Int64Attribute{
							Description: "Frequency to be used for synchronizing the integration.",
							Computed:    true,
						},
						"endpoint": schema.StringAttribute{
							Description: "Endpoint to be used for performing the integration.",
							Computed:    true,
						},
						"oauth_client_credentials": schema.SingleNestedAttribute{
							Description: "OAuth Client Credentials authentication information for the integration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"token_url": schema.StringAttribute{Required: true},
								"client_id": schema.StringAttribute{Required: true},
								"client_secret": schema.StringAttribute{
									Computed:  true,
									Sensitive: true,
								},
								"audience": schema.StringAttribute{Optional: true},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *integrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state integrationsDataSourceModel

	integrations, err := d.client.GetIntegrations(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Integrations",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, integration := range integrations {
		integrationState := convertIntegrationDTOToModel(ctx, integration, integrationResourceModel{})
		state.Integrations = append(state.Integrations, integrationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
