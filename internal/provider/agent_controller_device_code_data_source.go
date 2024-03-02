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
	_ datasource.DataSource              = &agentControllerDeviceCodeDataSource{}
	_ datasource.DataSourceWithConfigure = &agentControllerDeviceCodeDataSource{}
)

// NewagentControllerDeviceCodeDataSource is a helper function to simplify the provider implementation.
func NewAgentControllerDeviceCodeDataSource() datasource.DataSource {
	return &agentControllerDeviceCodeDataSource{}
}

// agentControllerDeviceCodeDataSource is the data source implementation.
type agentControllerDeviceCodeDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *agentControllerDeviceCodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *agentControllerDeviceCodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_controller_device_code"
}

// Schema defines the schema for the resource.
func (d *agentControllerDeviceCodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates an agent controller device code.",
		Attributes: map[string]schema.Attribute{
			"agent_controller_id": schema.StringAttribute{
				Description: "Unique identifier of the Agent Controller.",
				Required:    true,
			},
			"device_code": schema.StringAttribute{
				Description: "Generated Device Code of the Agent Controller.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *agentControllerDeviceCodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var deviceCodeRequest agentControllerDeviceCodeDataSourceModel
	var state agentControllerDeviceCodeDataSourceModel

	// Retrieve Agent Controller ID from plan
	resp.Diagnostics.Append(req.Config.Get(ctx, &deviceCodeRequest)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agentControllerID := deviceCodeRequest.ID.ValueString()

	deviceCode, err := d.client.GetAgentControllerDeviceCode(agentControllerID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Retrieve Aembit Agent Controller Device Code",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.DeviceCode = types.StringValue(deviceCode.DeviceCode)
	state.ID = types.StringValue(agentControllerID)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
