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
	_ datasource.DataSource              = &trustProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &trustProvidersDataSource{}
)

// NewTrustProvidersDataSource is a helper function to simplify the provider implementation.
func NewTrustProvidersDataSource() datasource.DataSource {
	return &trustProvidersDataSource{}
}

// trustProvidersDataSource is the data source implementation.
type trustProvidersDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *trustProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Metadata returns the data source type name.
func (d *trustProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_providers"
}

// Schema defines the schema for the resource.
func (d *trustProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an trust provider.",
		Attributes: map[string]schema.Attribute{
			"trust_providers": schema.ListNestedAttribute{
				Description: "List of trust providers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the trust provider.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the trust provider.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the trust provider.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the trust provider.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *trustProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state trustProvidersDataSourceModel

	trustProviders, err := d.client.GetTrustProviders(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Trust Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, trustProvider := range trustProviders {
		trustProviderState := trustProviderResourceModel{
			ID:          types.StringValue(trustProvider.EntityDTO.ExternalID),
			Name:        types.StringValue(trustProvider.EntityDTO.Name),
			Description: types.StringValue(trustProvider.EntityDTO.Description),
			IsActive:    types.BoolValue(trustProvider.EntityDTO.IsActive),
		}
		state.TrustProviders = append(state.TrustProviders, trustProviderState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
