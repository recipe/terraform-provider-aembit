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
			fmt.Sprintf("Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
							Description: "Unique identifier of the trust provider.",
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
						"tags": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"azure_metadata": schema.SingleNestedAttribute{
							Description: "Azure Metadata type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"sku":   schema.StringAttribute{Computed: true},
								"vm_id": schema.StringAttribute{Computed: true},
								"subscription_id": schema.StringAttribute{
									Computed: true,
									//Validators: []validator.String{
									//	// Validate azure_metadata has at least one value
									//	stringvalidator.AtLeastOneOf(path.Expressions{
									//		path.MatchRelative().AtParent().AtName("sku"),
									//		path.MatchRelative().AtParent().AtName("vm_id"),
									//	}...),
									//},
								},
							},
						},
						"aws_metadata": schema.SingleNestedAttribute{
							Description: "AWS Metadata type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"certificate": schema.StringAttribute{
									Computed:    true,
									Description: "PEM Certificate to be used for Signature verification",
								},
								"account_id":                schema.StringAttribute{Computed: true},
								"architecture":              schema.StringAttribute{Computed: true},
								"availability_zone":         schema.StringAttribute{Computed: true},
								"billing_products":          schema.StringAttribute{Computed: true},
								"image_id":                  schema.StringAttribute{Computed: true},
								"instance_id":               schema.StringAttribute{Computed: true},
								"instance_type":             schema.StringAttribute{Computed: true},
								"kernel_id":                 schema.StringAttribute{Computed: true},
								"marketplace_product_codes": schema.StringAttribute{Computed: true},
								"pending_time":              schema.StringAttribute{Computed: true},
								"private_ip":                schema.StringAttribute{Computed: true},
								"ramdisk_id":                schema.StringAttribute{Computed: true},
								"region":                    schema.StringAttribute{Computed: true},
								"version":                   schema.StringAttribute{Computed: true},
							},
						},
						"kerberos": schema.SingleNestedAttribute{
							Description: "Kerberos type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"agent_controller_ids": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"principal": schema.StringAttribute{Computed: true},
								"realm":     schema.StringAttribute{Computed: true},
								"source_ip": schema.StringAttribute{Computed: true},
							},
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
		trustProviderState := convertTrustProviderDTOToModel(ctx, trustProvider)
		state.TrustProviders = append(state.TrustProviders, trustProviderState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
