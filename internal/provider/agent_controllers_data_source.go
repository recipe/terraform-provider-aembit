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
	_ datasource.DataSource              = &agentControllersDataSource{}
	_ datasource.DataSourceWithConfigure = &agentControllersDataSource{}
)

// NewAgentControllersDataSource is a helper function to simplify the provider implementation.
func NewAgentControllersDataSource() datasource.DataSource {
	return &agentControllersDataSource{}
}

// agentControllersDataSource is the data source implementation.
type agentControllersDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *agentControllersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *agentControllersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_controllers"
}

// Schema defines the schema for the resource.
func (d *agentControllersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an agent controller.",
		Attributes: map[string]schema.Attribute{
			"agent_controllers": schema.ListNestedAttribute{
				Description: "List of agent controllers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the agent controller.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the agent controller.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the agent controller.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the agent controller.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Tags are key-value pairs.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"trust_provider_id": schema.BoolAttribute{
							Description: "Trust Provider to use for authentication of the agent controller.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *agentControllersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state agentControllersDataSourceModel

	agentControllers, err := d.client.GetAgentControllers(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Agent Controllers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, agentController := range agentControllers {
		agentControllerState := ConvertAgentControllerDTOToModel(ctx, agentController)
		state.AgentControllers = append(state.AgentControllers, agentControllerState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
