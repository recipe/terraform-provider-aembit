package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &trustProviderResource{}
	_ resource.ResourceWithConfigure   = &trustProviderResource{}
	_ resource.ResourceWithImportState = &trustProviderResource{}
)

// NewTrustProviderResource is a helper function to simplify the provider implementation.
func NewTrustProviderResource() resource.Resource {
	return &trustProviderResource{}
}

// trustProviderResource is the resource implementation.
type trustProviderResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *trustProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_provider"
}

// Configure adds the provider configured client to the resource.
func (r *trustProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Schema defines the schema for the resource.
func (r *trustProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Trust Provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Trust Provider.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for the Trust Provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Trust Provider.",
				Optional:    true,
				Computed:    true,
			},
			"azure_metadata": schema.SingleNestedAttribute{
				Description: "Azure Metadata type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sku": schema.StringAttribute{
						Description: "Specific SKU for the Virtual Machine image.",
						Optional:    true,
					},
					"vm_id": schema.StringAttribute{
						Description: "Unique identifier for the Virtual Machine.",
						Optional:    true,
					},
					"subscription_id": schema.StringAttribute{
						Description: "Azure subscription for the Virtual Machine.",
						Optional:    true,
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
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"certificate": schema.StringAttribute{
						Description: "PEM Certificate to be used for Signature verification.",
						Optional:    true,
					},
					"account_id": schema.StringAttribute{
						Description: "The ID of the AWS account that launched the instance.",
						Optional:    true,
					},
					"architecture": schema.StringAttribute{
						Description: "The architecture of the AMI used to launch the instance (i386 | x86_64 | arm64).",
						Optional:    true,
					},
					"availability_zone": schema.StringAttribute{
						Description: "The Availability Zone in which the instance is running.",
						Optional:    true,
					},
					"billing_products": schema.StringAttribute{
						Description: "The billing products of the instance.",
						Optional:    true,
					},
					"image_id": schema.StringAttribute{
						Description: "The ID of the AMI used to launch the instance.",
						Optional:    true,
					},
					"instance_id": schema.StringAttribute{
						Description: "The ID of the instance.",
						Optional:    true,
					},
					"instance_type": schema.StringAttribute{
						Description: "The instance type of the instance.",
						Optional:    true,
					},
					"kernel_id": schema.StringAttribute{
						Description: "The ID of the kernel associated with the instance, if applicable.",
						Optional:    true,
					},
					"marketplace_product_codes": schema.StringAttribute{
						Description: "The AWS Marketplace product code of the AMI used to launch the instance.",
						Optional:    true,
					},
					"pending_time": schema.StringAttribute{
						Description: "The date and time that the instance was launched.",
						Optional:    true,
					},
					"private_ip": schema.StringAttribute{
						Description: "The private IPv4 address of the instance.",
						Optional:    true,
					},
					"ramdisk_id": schema.StringAttribute{
						Description: "The ID of the RAM disk associated with the instance, if applicable.",
						Optional:    true,
					},
					"region": schema.StringAttribute{
						Description: "The Region in which the instance is running.",
						Optional:    true,
					},
					"version": schema.StringAttribute{
						Description: "The version of the instance identity document format.",
						Optional:    true,
					},
				},
			},
			"kerberos": schema.SingleNestedAttribute{
				Description: "Kerberos type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"agent_controller_id": schema.StringAttribute{
						Description: "Unique identifier for the Aembit Agent Controller to use for Signature verification.",
						Optional:    true,
					},
					"principal": schema.StringAttribute{
						Description: "The Kerberos Principal of the authenticated Agent Proxy.",
						Optional:    true,
					},
					"realm": schema.StringAttribute{
						Description: "The Kerberos Realm of the authenticated Agent Proxy.",
						Optional:    true,
					},
					"source_ip": schema.StringAttribute{
						Description: "The Source IP Address of the authenticated Agent Proxy.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one Trust Provider type is specified.
func (r *trustProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("azure_metadata"),
			path.MatchRoot("aws_metadata"),
			path.MatchRoot("kerberos"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *trustProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(plan, nil)

	// Create new Trust Provider
	trustProvider, err := r.client.CreateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trust Provider",
			"Could not create Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertTrustProviderDTOToModel(*trustProvider)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *trustProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	trustProvider, err := r.client.GetTrustProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Trust Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertTrustProviderDTOToModel(trustProvider)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trustProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(plan, &externalID)

	// Update Trust Provider
	trustProvider, err := r.client.UpdateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Trust Provider",
			"Could not update Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertTrustProviderDTOToModel(*trustProvider)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *trustProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Trust Provider is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Trust Provider",
			"Trust Provider is active and cannot be deleted. Please mark the trust as inactive first.",
		)
		return
	}

	// Delete existing Trust Provider
	_, err := r.client.DeleteTrustProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Trust Provider",
			"Could not delete Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *trustProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertTrustProviderModelToDTO(model trustProviderResourceModel, externalID *string) aembit.TrustProviderDTO {
	var trust aembit.TrustProviderDTO
	trust.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if externalID != nil {
		trust.EntityDTO.ExternalID = *externalID
	}

	// Handle the Azure Metadata use case
	if model.AzureMetadata != nil {
		convertAzureMetadataModelToDTO(model, &trust)
	}
	if model.AwsMetadata != nil {
		convertAwsMetadataModelToDTO(model, &trust)
	}
	if model.Kerberos != nil {
		convertKerberosModelToDTO(model, &trust)
	}

	return trust
}

func convertAzureMetadataModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AzureMetadataService"
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)

	if len(model.AzureMetadata.Sku.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AzureSku", Value: model.AzureMetadata.Sku.ValueString(),
		})
	}
	if len(model.AzureMetadata.VMID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AzureVmId", Value: model.AzureMetadata.VMID.ValueString(),
		})
	}
	if len(model.AzureMetadata.SubscriptionID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AzureSubscriptionId", Value: model.AzureMetadata.SubscriptionID.ValueString(),
		})
	}
}

func convertAwsMetadataModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AWSMetadataService"
	dto.Certificate = base64.StdEncoding.EncodeToString([]byte(model.AwsMetadata.Certificate.ValueString()))
	dto.PemType = "Certificate"
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)

	if len(model.AwsMetadata.AccountID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsAccountId", Value: model.AwsMetadata.AccountID.ValueString(),
		})
	}
	if len(model.AwsMetadata.Architecture.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsArchitecture", Value: model.AwsMetadata.Architecture.ValueString(),
		})
	}
	if len(model.AwsMetadata.AvailabilityZone.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsAvailabilityZone", Value: model.AwsMetadata.AvailabilityZone.ValueString(),
		})
	}
	if len(model.AwsMetadata.BillingProducts.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsBillingProducts", Value: model.AwsMetadata.BillingProducts.ValueString(),
		})
	}
	if len(model.AwsMetadata.ImageID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsImageId", Value: model.AwsMetadata.ImageID.ValueString(),
		})
	}
	if len(model.AwsMetadata.InstanceID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsInstanceId", Value: model.AwsMetadata.InstanceID.ValueString(),
		})
	}
	if len(model.AwsMetadata.InstanceType.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsInstanceType", Value: model.AwsMetadata.InstanceType.ValueString(),
		})
	}
	if len(model.AwsMetadata.KernelID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsKernelId", Value: model.AwsMetadata.KernelID.ValueString(),
		})
	}
	if len(model.AwsMetadata.MarketplaceProductCodes.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsMarketplaceProductCodes", Value: model.AwsMetadata.MarketplaceProductCodes.ValueString(),
		})
	}
	if len(model.AwsMetadata.PendingTime.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsPendingTime", Value: model.AwsMetadata.PendingTime.ValueString(),
		})
	}
	if len(model.AwsMetadata.PrivateIP.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsPrivateIp", Value: model.AwsMetadata.PrivateIP.ValueString(),
		})
	}
	if len(model.AwsMetadata.RamdiskID.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsRamdiskId", Value: model.AwsMetadata.RamdiskID.ValueString(),
		})
	}
	if len(model.AwsMetadata.Region.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsRegion", Value: model.AwsMetadata.Region.ValueString(),
		})
	}
	if len(model.AwsMetadata.Version.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsVersion", Value: model.AwsMetadata.Version.ValueString(),
		})
	}
}

func convertKerberosModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "Kerberos"
	dto.AgentControllerID = model.Kerberos.AgentControllerID.ValueString()
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)

	if len(model.Kerberos.Principal.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "Principal", Value: model.Kerberos.Principal.ValueString(),
		})
	}
	if len(model.Kerberos.Realm.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "Realm", Value: model.Kerberos.Realm.ValueString(),
		})
	}
	if len(model.Kerberos.SourceIP.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "SourceIp", Value: model.Kerberos.SourceIP.ValueString(),
		})
	}
}

func convertTrustProviderDTOToModel(dto aembit.TrustProviderDTO) trustProviderResourceModel {
	var model trustProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	switch dto.Provider {
	case "AzureMetadataService": // Azure Metadata
		model.AzureMetadata = convertAzureMetadataDTOToModel(dto)
	case "AWSMetadataService": // AWS Metadata
		model.AwsMetadata = convertAwsMetadataDTOToModel(dto)
	case "Kerberos": // Kerberos
		model.Kerberos = convertKerberosDTOToModel(dto)
	}

	return model
}

func convertAzureMetadataDTOToModel(dto aembit.TrustProviderDTO) *trustProviderAzureMetadataModel {
	model := &trustProviderAzureMetadataModel{
		Sku:            types.StringNull(),
		VMID:           types.StringNull(),
		SubscriptionID: types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AzureSku":
			model.Sku = types.StringValue(rule.Value)
		case "AzureVmId":
			model.VMID = types.StringValue(rule.Value)
		case "AzureSubscriptionId":
			model.SubscriptionID = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertAwsMetadataDTOToModel(dto aembit.TrustProviderDTO) *trustProviderAwsMetadataModel {
	decodedCert, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &trustProviderAwsMetadataModel{
		Certificate:             types.StringValue(string(decodedCert)),
		AccountID:               types.StringNull(),
		Architecture:            types.StringNull(),
		AvailabilityZone:        types.StringNull(),
		BillingProducts:         types.StringNull(),
		ImageID:                 types.StringNull(),
		InstanceID:              types.StringNull(),
		InstanceType:            types.StringNull(),
		KernelID:                types.StringNull(),
		MarketplaceProductCodes: types.StringNull(),
		PendingTime:             types.StringNull(),
		PrivateIP:               types.StringNull(),
		RamdiskID:               types.StringNull(),
		Region:                  types.StringNull(),
		Version:                 types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AwsAccountId":
			model.AccountID = types.StringValue(rule.Value)
		case "AwsArchitecture":
			model.Architecture = types.StringValue(rule.Value)
		case "AwsAvailabilityZone":
			model.AvailabilityZone = types.StringValue(rule.Value)
		case "AwsBillingProducts":
			model.BillingProducts = types.StringValue(rule.Value)
		case "AwsImageId":
			model.ImageID = types.StringValue(rule.Value)
		case "AwsInstanceId":
			model.InstanceID = types.StringValue(rule.Value)
		case "AwsInstanceType":
			model.InstanceType = types.StringValue(rule.Value)
		case "AwsKernelId":
			model.KernelID = types.StringValue(rule.Value)
		case "AwsMarketplaceProductCodes":
			model.MarketplaceProductCodes = types.StringValue(rule.Value)
		case "AwsPendingTime":
			model.PendingTime = types.StringValue(rule.Value)
		case "AwsPrivateIp":
			model.PrivateIP = types.StringValue(rule.Value)
		case "AwsRamdiskId":
			model.RamdiskID = types.StringValue(rule.Value)
		case "AwsRegion":
			model.Region = types.StringValue(rule.Value)
		case "AwsVersion":
			model.Version = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertKerberosDTOToModel(dto aembit.TrustProviderDTO) *trustProviderKerberosModel {
	model := &trustProviderKerberosModel{
		AgentControllerID: types.StringValue(dto.AgentControllerID),
		Principal:         types.StringNull(),
		Realm:             types.StringNull(),
		SourceIP:          types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "Principal":
			model.Principal = types.StringValue(rule.Value)
		case "Realm":
			model.Realm = types.StringValue(rule.Value)
		case "SourceIp":
			model.SourceIP = types.StringValue(rule.Value)
		}
	}
	return model
}
