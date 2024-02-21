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
	client *aembit.AembitClient
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

	client, ok := req.ProviderData.(*aembit.AembitClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.AembitClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
				Description: "Alphanumeric identifier of the trust provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the trust provider.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the trust provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the trust provider.",
				Optional:    true,
				Computed:    true,
			},
			"azure_metadata": schema.SingleNestedAttribute{
				Description: "Azure Metadata type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sku":   schema.StringAttribute{Optional: true},
					"vm_id": schema.StringAttribute{Optional: true},
					"subscription_id": schema.StringAttribute{
						Optional: true,
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
						Optional:    true,
						Description: "PEM Certificate to be used for Signature verification",
					},
					"account_id":                schema.StringAttribute{Optional: true},
					"architecture":              schema.StringAttribute{Optional: true},
					"availability_zone":         schema.StringAttribute{Optional: true},
					"billing_products":          schema.StringAttribute{Optional: true},
					"image_id":                  schema.StringAttribute{Optional: true},
					"instance_id":               schema.StringAttribute{Optional: true},
					"instance_type":             schema.StringAttribute{Optional: true},
					"kernel_id":                 schema.StringAttribute{Optional: true},
					"marketplace_product_codes": schema.StringAttribute{Optional: true},
					"pending_time":              schema.StringAttribute{Optional: true},
					"private_ip":                schema.StringAttribute{Optional: true},
					"ramdisk_id":                schema.StringAttribute{Optional: true},
					"region":                    schema.StringAttribute{Optional: true},
					"version":                   schema.StringAttribute{Optional: true},
				},
			},
			"kerberos": schema.SingleNestedAttribute{
				Description: "Kerberos type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"agent_controller_id": schema.StringAttribute{Optional: true},
					"principal":           schema.StringAttribute{Optional: true},
					"realm":               schema.StringAttribute{Optional: true},
					"source_ip":           schema.StringAttribute{Optional: true},
				},
			},
		},
	}
}

// Configure validators to ensure that only one trust provider type is specified
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
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, nil)

	// Create new Trust Provider
	trust_provider, err := r.client.CreateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust provider",
			"Could not create trust provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertTrustProviderDTOToModel(ctx, *trust_provider)

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
	trust_provider, err := r.client.GetTrustProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Trust Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertTrustProviderDTOToModel(ctx, trust_provider)

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
	external_id := state.ID.ValueString()

	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, &external_id)

	// Update Trust Provider
	trust_provider, err := r.client.UpdateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating trust provider",
			"Could not update trust provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertTrustProviderDTOToModel(ctx, *trust_provider)

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
			"Could not delete trust provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId
func (r *trustProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertTrustProviderModelToDTO(ctx context.Context, model trustProviderResourceModel, external_id *string) aembit.TrustProviderDTO {
	var trust aembit.TrustProviderDTO
	trust.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if external_id != nil {
		trust.EntityDTO.ExternalId = *external_id
	}

	// Handle the Azure Metadata use case
	if model.AzureMetadata != nil {
		convertAzureMetadataModelToDTO(ctx, model, &trust)
	}
	if model.AwsMetadata != nil {
		convertAwsMetadataModelToDTO(ctx, model, &trust)
	}
	if model.Kerberos != nil {
		convertKerberosModelToDTO(ctx, model, &trust)
	}

	return trust
}

func convertAzureMetadataModelToDTO(ctx context.Context, model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AzureMetadataService"
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)

	if len(model.AzureMetadata.Sku.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AzureSku", Value: model.AzureMetadata.Sku.ValueString(),
		})
	}
	if len(model.AzureMetadata.VmId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AzureVmId", Value: model.AzureMetadata.VmId.ValueString(),
		})
	}
	if len(model.AzureMetadata.SubscriptionId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AzureSubscriptionId", Value: model.AzureMetadata.SubscriptionId.ValueString(),
		})
	}
}

func convertAwsMetadataModelToDTO(ctx context.Context, model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AWSMetadataService"
	dto.Certificate = base64.StdEncoding.EncodeToString([]byte(model.AwsMetadata.Certificate.ValueString()))
	dto.PemType = "Certificate"
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)

	if len(model.AwsMetadata.AccountId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsAccountId", Value: model.AwsMetadata.AccountId.ValueString(),
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
	if len(model.AwsMetadata.ImageId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsImageId", Value: model.AwsMetadata.ImageId.ValueString(),
		})
	}
	if len(model.AwsMetadata.InstanceId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsInstanceId", Value: model.AwsMetadata.InstanceId.ValueString(),
		})
	}
	if len(model.AwsMetadata.InstanceType.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsInstanceType", Value: model.AwsMetadata.InstanceType.ValueString(),
		})
	}
	if len(model.AwsMetadata.KernelId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsKernelId", Value: model.AwsMetadata.KernelId.ValueString(),
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
	if len(model.AwsMetadata.RamdiskId.ValueString()) > 0 {
		dto.MatchRules = append(dto.MatchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: "AwsRamdiskId", Value: model.AwsMetadata.RamdiskId.ValueString(),
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

func convertKerberosModelToDTO(ctx context.Context, model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "Kerberos"
	dto.AgentControllerId = model.Kerberos.AgentControllerId.ValueString()
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

func convertTrustProviderDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO) trustProviderResourceModel {
	var model trustProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalId)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	switch dto.Provider {
	case "AzureMetadataService": // Azure Metadata
		model.AzureMetadata = convertAzureMetadataDTOToModel(ctx, dto)
	case "AWSMetadataService": // AWS Metadata
		model.AwsMetadata = convertAwsMetadataDTOToModel(ctx, dto)
	case "Kerberos": // Kerberos
		model.Kerberos = convertKerberosDTOToModel(ctx, dto)
	}

	return model
}

func convertAzureMetadataDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO) *trustProviderAzureMetadataModel {
	model := &trustProviderAzureMetadataModel{
		Sku:            types.StringNull(),
		VmId:           types.StringNull(),
		SubscriptionId: types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AzureSku":
			model.Sku = types.StringValue(rule.Value)
		case "AzureVmId":
			model.VmId = types.StringValue(rule.Value)
		case "AzureSubscriptionId":
			model.SubscriptionId = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertAwsMetadataDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO) *trustProviderAwsMetadataModel {
	decodedCert, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &trustProviderAwsMetadataModel{
		Certificate:             types.StringValue(string(decodedCert)),
		AccountId:               types.StringNull(),
		Architecture:            types.StringNull(),
		AvailabilityZone:        types.StringNull(),
		BillingProducts:         types.StringNull(),
		ImageId:                 types.StringNull(),
		InstanceId:              types.StringNull(),
		InstanceType:            types.StringNull(),
		KernelId:                types.StringNull(),
		MarketplaceProductCodes: types.StringNull(),
		PendingTime:             types.StringNull(),
		PrivateIP:               types.StringNull(),
		RamdiskId:               types.StringNull(),
		Region:                  types.StringNull(),
		Version:                 types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AwsAccountId":
			model.AccountId = types.StringValue(rule.Value)
		case "AwsArchitecture":
			model.Architecture = types.StringValue(rule.Value)
		case "AwsAvailabilityZone":
			model.AvailabilityZone = types.StringValue(rule.Value)
		case "AwsBillingProducts":
			model.BillingProducts = types.StringValue(rule.Value)
		case "AwsImageId":
			model.ImageId = types.StringValue(rule.Value)
		case "AwsInstanceId":
			model.InstanceId = types.StringValue(rule.Value)
		case "AwsInstanceType":
			model.InstanceType = types.StringValue(rule.Value)
		case "AwsKernelId":
			model.KernelId = types.StringValue(rule.Value)
		case "AwsMarketplaceProductCodes":
			model.MarketplaceProductCodes = types.StringValue(rule.Value)
		case "AwsPendingTime":
			model.PendingTime = types.StringValue(rule.Value)
		case "AwsPrivateIp":
			model.PrivateIP = types.StringValue(rule.Value)
		case "AwsRamdiskId":
			model.RamdiskId = types.StringValue(rule.Value)
		case "AwsRegion":
			model.Region = types.StringValue(rule.Value)
		case "AwsVersion":
			model.Version = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertKerberosDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO) *trustProviderKerberosModel {
	model := &trustProviderKerberosModel{
		AgentControllerId: types.StringValue(dto.AgentControllerId),
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
