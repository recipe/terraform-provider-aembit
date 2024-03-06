package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// trustProviderResourceModel maps the resource schema.
type trustProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                 types.String                     `tfsdk:"id"`
	Name               types.String                     `tfsdk:"name"`
	Description        types.String                     `tfsdk:"description"`
	IsActive           types.Bool                       `tfsdk:"is_active"`
	Tags               types.Map                        `tfsdk:"tags"`
	AzureMetadata      *trustProviderAzureMetadataModel `tfsdk:"azure_metadata"`
	AwsEcsRole         *trustProviderAwsEcsRoleModel    `tfsdk:"aws_ecs_role"`
	AwsMetadata        *trustProviderAwsMetadataModel   `tfsdk:"aws_metadata"`
	GcpIdentity        *trustProviderGcpIdentityModel   `tfsdk:"gcp_identity"`
	GitHubAction       *trustProviderGitHubActionModel  `tfsdk:"github_action"`
	Kerberos           *trustProviderKerberosModel      `tfsdk:"kerberos"`
	KubernetesService  *trustProviderKubernetesModel    `tfsdk:"kubernetes_service_account"`
	TerraformWorkspace *trustProviderTerraformModel     `tfsdk:"terraform_workspace"`
}

// trustProviderDataSourceModel maps the datasource schema.
type trustProvidersDataSourceModel struct {
	TrustProviders []trustProviderResourceModel `tfsdk:"trust_providers"`
}

type trustProviderAzureMetadataModel struct {
	Sku            types.String `tfsdk:"sku"`
	VMID           types.String `tfsdk:"vm_id"`
	SubscriptionID types.String `tfsdk:"subscription_id"`
}

type trustProviderAwsEcsRoleModel struct {
	AccountID   types.String `tfsdk:"account_id"`
	AssumedRole types.String `tfsdk:"assumed_role"`
	RoleARN     types.String `tfsdk:"role_arn"`
	Username    types.String `tfsdk:"username"`
}

type trustProviderAwsMetadataModel struct {
	Certificate             types.String `tfsdk:"certificate"`
	AccountID               types.String `tfsdk:"account_id"`
	Architecture            types.String `tfsdk:"architecture"`
	AvailabilityZone        types.String `tfsdk:"availability_zone"`
	BillingProducts         types.String `tfsdk:"billing_products"`
	ImageID                 types.String `tfsdk:"image_id"`
	InstanceID              types.String `tfsdk:"instance_id"`
	InstanceType            types.String `tfsdk:"instance_type"`
	KernelID                types.String `tfsdk:"kernel_id"`
	MarketplaceProductCodes types.String `tfsdk:"marketplace_product_codes"`
	PendingTime             types.String `tfsdk:"pending_time"`
	PrivateIP               types.String `tfsdk:"private_ip"`
	RamdiskID               types.String `tfsdk:"ramdisk_id"`
	Region                  types.String `tfsdk:"region"`
	Version                 types.String `tfsdk:"version"`
}

type trustProviderKerberosModel struct {
	AgentControllerIDs []types.String `tfsdk:"agent_controller_ids"`
	Principal          types.String   `tfsdk:"principal"`
	Realm              types.String   `tfsdk:"realm"`
	SourceIP           types.String   `tfsdk:"source_ip"`
}

type trustProviderKubernetesModel struct {
	Issuer             types.String `tfsdk:"issuer"`
	Namespace          types.String `tfsdk:"namespace"`
	PodName            types.String `tfsdk:"pod_name"`
	ServiceAccountName types.String `tfsdk:"service_account_name"`
	Subject            types.String `tfsdk:"subject"`
	OIDCEndpoint       types.String `tfsdk:"oidc_endpoint"`
	PublicKey          types.String `tfsdk:"public_key"`
}

type trustProviderGcpIdentityModel struct {
	EMail types.String `tfsdk:"email"`
}

type trustProviderGitHubActionModel struct {
	Actor      types.String `tfsdk:"actor"`
	Repository types.String `tfsdk:"repository"`
	Workflow   types.String `tfsdk:"workflow"`
}

type trustProviderTerraformModel struct {
	OrganizationID types.String `tfsdk:"organization_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	WorkspaceID    types.String `tfsdk:"workspace_id"`
}
