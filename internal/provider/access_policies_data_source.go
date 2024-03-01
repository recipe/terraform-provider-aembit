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
	_ datasource.DataSource              = &accessPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &accessPoliciesDataSource{}
)

// NewAccessPoliciesDataSource is a helper function to simplify the provider implementation.
func NewAccessPoliciesDataSource() datasource.DataSource {
	return &accessPoliciesDataSource{}
}

// accessPoliciesDataSource is the data source implementation.
type accessPoliciesDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *accessPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *accessPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_policies"
}

// Schema defines the schema for the resource.
func (d *accessPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages access policies.",
		Attributes: map[string]schema.Attribute{
			"access_policies": schema.ListNestedAttribute{
				Description: "List of access policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the access policy.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the access policy.",
							Computed:    true,
						},
						"client_workload": schema.StringAttribute{
							Description: "Configured client workload of the access policy.",
							Computed:    true,
						},
						"trust_providers": schema.SetAttribute{
							Description: "Set of Trust Providers to enforce on the Access Policy.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"access_conditions": schema.SetAttribute{
							Description: "Set of Access Conditions to enforce on the Access Policy.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"credential_provider": schema.StringAttribute{
							Description: "Credential Provider ID configured in the Access Policy.",
							Computed:    true,
						},
						"server_workload": schema.StringAttribute{
							Description: "Configured server workload of the access policy.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *accessPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accessPoliciesDataSourceModel

	accessPolicies, err := d.client.GetAccessPolicies(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Trust Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, accessPolicy := range accessPolicies {
		accessPolicyState := ConvertAccessPolicyExternalDTOToModel(accessPolicy)
		state.AccessPolicies = append(state.AccessPolicies, accessPolicyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
