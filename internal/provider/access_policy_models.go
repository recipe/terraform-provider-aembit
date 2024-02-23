package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// accessPolicyResourceModel maps the resource schema.
type accessPolicyResourceModel struct {
	// ID is required for Framework acceptance testing
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	IsActive       types.Bool   `tfsdk:"is_active"`
	ClientWorkload types.String `tfsdk:"client_workload"`
	ServerWorkload types.String `tfsdk:"server_workload"`
	//CredentialProvider types.String      `tfsdk:"credential_provider"`
	//TrustProviders     []types.String    `tfsdk:"trust_providers"`
	//PolicyNotes        []policyNoteModel `tfsdk:"policy_notes"`
}

// accessPoliciesDataSourceModel maps the datasource schema.
type accessPoliciesDataSourceModel struct {
	AccessPolicies []accessPolicyResourceModel `tfsdk:"access_policies"`
}

// policyNoteModel maps the datasource schema.
//type policyNoteModel struct {
//	Note types.String `tfsdk:"note"`
//}
