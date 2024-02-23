package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// trustProviderResourceModel maps the resource schema.
type trustProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID            types.String                     `tfsdk:"id"`
	Name          types.String                     `tfsdk:"name"`
	Description   types.String                     `tfsdk:"description"`
	IsActive      types.Bool                       `tfsdk:"is_active"`
	AzureMetadata *trustProviderAzureMetadataModel `tfsdk:"azure_metadata"`
	AwsMetadata   *trustProviderAwsMetadataModel   `tfsdk:"aws_metadata"`
	Kerberos      *trustProviderKerberosModel      `tfsdk:"kerberos"`
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
	AgentControllerID types.String `tfsdk:"agent_controller_id"`
	Principal         types.String `tfsdk:"principal"`
	Realm             types.String `tfsdk:"realm"`
	SourceIP          types.String `tfsdk:"source_ip"`
}
