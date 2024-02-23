package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// clientWorkloadResourceModel maps the resource schema.
type clientWorkloadResourceModel struct {
	// ID is required for Framework acceptance testing
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	IsActive    types.Bool        `tfsdk:"is_active"`
	Identities  []identitiesModel `tfsdk:"identities"`
	Type        types.String      `tfsdk:"type"`
}

// clientWorkloadDataSourceModel maps the datasource schema.
//type clientWorkloadsDataSourceModel struct {
//	ClientWorkloads []clientWorkloadResourceModel `tfsdk:"client_workloads"`
//}

// identitiesModel maps client workload identity data.
type identitiesModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}
