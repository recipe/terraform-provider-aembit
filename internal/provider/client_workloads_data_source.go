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
	_ datasource.DataSource              = &clientWorkloadsDataSource{}
	_ datasource.DataSourceWithConfigure = &clientWorkloadsDataSource{}
)

// NewClientWorkloadsDataSource is a helper function to simplify the provider implementation.
func NewClientWorkloadsDataSource() datasource.DataSource {
	return &clientWorkloadsDataSource{}
}

// clientWorkloadsDataSource is the data source implementation.
type clientWorkloadsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *clientWorkloadsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *clientWorkloadsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_workloads"
}

// Schema defines the schema for the resource.
func (r *clientWorkloadsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages client workloads.",
		Attributes: map[string]schema.Attribute{
			"client_workloads": schema.ListNestedAttribute{
				Description: "List of client workloads.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the client workload.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the client workload.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the client workload.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the client workload.",
							Computed:    true,
						},
						"identities": schema.SetNestedAttribute{
							Description: "Set of client workload identities.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Description: "Client identity type.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Client identity value.",
										Computed:    true,
									},
								},
							},
						},
						"tags": schema.MapAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *clientWorkloadsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clientWorkloadsDataSourceModel

	clientWorkloads, err := d.client.GetClientWorkloads(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Trust Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, clientWorkload := range clientWorkloads {
		clientWorkloadState := convertClientWorkloadDTOToModel(ctx, clientWorkload)
		state.ClientWorkloads = append(state.ClientWorkloads, clientWorkloadState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
