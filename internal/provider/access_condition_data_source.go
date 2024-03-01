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
	_ datasource.DataSource              = &accessConditionsDataSource{}
	_ datasource.DataSourceWithConfigure = &accessConditionsDataSource{}
)

// NewAccessConditionsDataSource is a helper function to simplify the provider implementation.
func NewAccessConditionsDataSource() datasource.DataSource {
	return &accessConditionsDataSource{}
}

// accessConditionsDataSource is the data source implementation.
type accessConditionsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *accessConditionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *accessConditionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_conditions"
}

// Schema defines the schema for the resource.
func (d *accessConditionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an accessCondition.",
		Attributes: map[string]schema.Attribute{
			"access_conditions": schema.ListNestedAttribute{
				Description: "List of accessConditions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the accessCondition.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the accessCondition.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the accessCondition.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the accessCondition.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Tags are key-value pairs.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"integration_id": schema.StringAttribute{
							Description: "ID of the Integration used by the Access Condition.",
							Computed:    true,
						},
						"wiz_conditions": schema.SingleNestedAttribute{
							Description: "Wiz Specific rules for the Access Condition.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"max_last_seen":               schema.Int64Attribute{Required: true},
								"container_cluster_connected": schema.BoolAttribute{Required: true},
							},
						},
						"crowdstrike_conditions": schema.SingleNestedAttribute{
							Description: "CrowdStrike Specific rules for the Access Condition.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"max_last_seen":       schema.Int64Attribute{Required: true},
								"match_hostname":      schema.BoolAttribute{Required: true},
								"match_serial_number": schema.BoolAttribute{Required: true},
								"prevent_rfm":         schema.BoolAttribute{Required: true},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *accessConditionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accessConditionsDataSourceModel

	accessConditions, err := d.client.GetAccessConditions(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit AccessConditions",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, accessCondition := range accessConditions {
		accessConditionState := convertAccessConditionDTOToModel(ctx, accessCondition, accessConditionResourceModel{})
		state.AccessConditions = append(state.AccessConditions, accessConditionState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
