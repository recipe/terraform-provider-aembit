package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
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
					},
				},
			},
			"aws_ecs_role": schema.SingleNestedAttribute{
				Description: "AWS ECS Role type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "The ID of the AWS account that is hosting the ECS Task.",
						Optional:    true,
					},
					"assumed_role": schema.StringAttribute{
						Description: "The Name of the AWS IAM Role which is running the ECS Task.",
						Optional:    true,
					},
					"role_arn": schema.StringAttribute{
						Description: "The ARN of the AWS IAM Role which is running the ECS Task.",
						Optional:    true,
					},
					"username": schema.StringAttribute{
						Description: "The UsernID of the AWS IAM Account which is running the ECS Task (not commonly used).",
						Optional:    true,
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
			"gcp_identity": schema.SingleNestedAttribute{
				Description: "GCP Identity type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"email": schema.StringAttribute{
						Description: "The Email of the GCP Service Account used by the associated GCP resource.",
						Optional:    true,
					},
				},
			},
			"github_action": schema.SingleNestedAttribute{
				Description: "GitHub Action type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"actor": schema.StringAttribute{
						Description: "The GitHub Actor which initiated the GitHub Action.",
						Optional:    true,
					},
					"repository": schema.StringAttribute{
						Description: "The GitHub Repository associated with the GitHub Action ID Token.",
						Optional:    true,
					},
					"workflow": schema.StringAttribute{
						Description: "The GitHub Workflow execution associated with the GitHub Action ID Token.",
						Optional:    true,
					},
				},
			},
			"kerberos": schema.SingleNestedAttribute{
				Description: "Kerberos type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"agent_controller_ids": schema.SetAttribute{
						Description: "Unique identifier for the Aembit Agent Controller to use for Signature verification.",
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
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
			"kubernetes_service_account": schema.SingleNestedAttribute{
				Description: "Kubernetes Service Account type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"issuer": schema.StringAttribute{
						Description: "The Issuer (`iss` claim) of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"namespace": schema.StringAttribute{
						Description: "The Namespace of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"pod_name": schema.StringAttribute{
						Description: "The Pod Name of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"service_account_name": schema.StringAttribute{
						Description: "The Service Account Name of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"subject": schema.StringAttribute{
						Description: "The Subject (`sub` claim) of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"oidc_endpoint": schema.StringAttribute{
						Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"public_key": schema.StringAttribute{
						Description: "The Public Key that can be used to verify the signature of the Kubernetes Service Account Token.",
						Optional:    true,
					},
				},
			},
			"terraform_workspace": schema.SingleNestedAttribute{
				Description: "Terraform Workspace type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"organization_id": schema.StringAttribute{
						Description: "The Organization ID of the calling Terraform Workspace.",
						Optional:    true,
					},
					"project_id": schema.StringAttribute{
						Description: "The Project ID of the calling Terraform Workspace.",
						Optional:    true,
					},
					"workspace_id": schema.StringAttribute{
						Description: "The Workspace ID of the calling Terraform Workspace.",
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
			path.MatchRoot("aws_ecs_role"),
			path.MatchRoot("aws_metadata"),
			path.MatchRoot("azure_metadata"),
			path.MatchRoot("gcp_identity"),
			path.MatchRoot("github_action"),
			path.MatchRoot("kerberos"),
			path.MatchRoot("kubernetes_service_account"),
			path.MatchRoot("terraform_workspace"),
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
	trustProvider, err := r.client.CreateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trust Provider",
			"Could not create Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertTrustProviderDTOToModel(ctx, *trustProvider)

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

	state = convertTrustProviderDTOToModel(ctx, trustProvider)

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
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, &externalID)

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
	state = convertTrustProviderDTOToModel(ctx, *trustProvider)

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

	// Check if Trust Provider is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableTrustProvider(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Trust Provider",
				"Could not disable Trust Provider, unexpected error: "+err.Error(),
			)
			return
		}
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

// Model to DTO conversion methods.
func convertTrustProviderModelToDTO(ctx context.Context, model trustProviderResourceModel, externalID *string) aembit.TrustProviderDTO {
	var trust aembit.TrustProviderDTO
	trust.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			trust.Tags = append(trust.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}
	if externalID != nil {
		trust.EntityDTO.ExternalID = *externalID
	}

	// Transform the various Trust Provider types
	if model.AwsMetadata != nil {
		convertAwsMetadataModelToDTO(model, &trust)
	}
	if model.AwsEcsRole != nil {
		convertAwsEcsRoleModelToDTO(model, &trust)
	}
	if model.AzureMetadata != nil {
		convertAzureMetadataModelToDTO(model, &trust)
	}
	if model.GcpIdentity != nil {
		convertGcpIdentityModelToDTO(model, &trust)
	}
	if model.GitHubAction != nil {
		convertGitHubActionModelToDTO(model, &trust)
	}
	if model.Kerberos != nil {
		convertKerberosModelToDTO(model, &trust)
	}
	if model.KubernetesService != nil {
		convertKubernetesModelToDTO(model, &trust)
	}
	if model.TerraformWorkspace != nil {
		convertTerraformModelToDTO(model, &trust)
	}

	return trust
}

func appendMatchRuleIfExists(matchRules []aembit.TrustProviderMatchRuleDTO, value basetypes.StringValue, attrName string) []aembit.TrustProviderMatchRuleDTO {
	if len(value.ValueString()) > 0 {
		return append(matchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: attrName, Value: value.ValueString(),
		})
	}
	return matchRules
}

func convertAzureMetadataModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AzureMetadataService"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.Sku, "AzureSku")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.VMID, "AzureVmId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.SubscriptionID, "AzureSubscriptionId")
}

func convertAwsEcsRoleModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AWSECSRole"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsEcsRole.AccountID, "AwsAccountId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsEcsRole.AssumedRole, "AwsAssumedRole")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsEcsRole.RoleARN, "AwsRoleARN")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsEcsRole.Username, "AwsUsername")
}

func convertAwsMetadataModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AWSMetadataService"
	dto.Certificate = base64.StdEncoding.EncodeToString([]byte(model.AwsMetadata.Certificate.ValueString()))
	dto.PemType = "Certificate"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.AccountID, "AwsAccountId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Architecture, "AwsArchitecture")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.AvailabilityZone, "AwsAvailabilityZone")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.BillingProducts, "AwsBillingProducts")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.ImageID, "AwsImageId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.InstanceID, "AwsInstanceId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.InstanceType, "AwsInstanceType")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.KernelID, "AwsKernelId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.MarketplaceProductCodes, "AwsMarketplaceProductCodes")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.PendingTime, "AwsPendingTime")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.PrivateIP, "AwsPrivateIp")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.RamdiskID, "AwsRamdiskId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Region, "AwsRegion")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Version, "AwsVersion")
}

func convertGcpIdentityModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "GcpIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GcpIdentity.EMail, "Email")
}

func convertGitHubActionModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "GitHubIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitHubAction.Actor, "GithubActor")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitHubAction.Repository, "GithubRepository")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitHubAction.Workflow, "GithubWorkflow")
}

func convertKerberosModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "Kerberos"
	dto.AgentControllerIDs = make([]string, len(model.Kerberos.AgentControllerIDs))
	for i, controllerID := range model.Kerberos.AgentControllerIDs {
		dto.AgentControllerIDs[i] = controllerID.ValueString()
	}

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.Principal, "Principal")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.Realm, "Realm")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.SourceIP, "SourceIp")
}

func convertKubernetesModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "KubernetesServiceAccount"
	dto.Certificate = base64.StdEncoding.EncodeToString([]byte(model.KubernetesService.PublicKey.ValueString()))
	if len(dto.Certificate) > 0 {
		dto.PemType = "PublicKey"
	}
	dto.OidcUrl = model.KubernetesService.OIDCEndpoint.ValueString()

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.Issuer, "KubernetesIss")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.Namespace, "KubernetesIoNamespace")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.PodName, "KubernetesIoPodName")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.ServiceAccountName, "KubernetesIoServiceAccountName")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.Subject, "KubernetesSub")
}

func convertTerraformModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "TerraformIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.TerraformWorkspace.OrganizationID, "TerraformOrganizationId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.TerraformWorkspace.ProjectID, "TerraformProjectId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.TerraformWorkspace.WorkspaceID, "TerraformWorkspaceId")
}

// DTO to Model conversion methods.
func convertTrustProviderDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO) trustProviderResourceModel {
	var model trustProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	switch dto.Provider {
	case "AWSECSRole":
		model.AwsEcsRole = convertAwsEcsRoleDTOToModel(dto)
	case "AWSMetadataService":
		model.AwsMetadata = convertAwsMetadataDTOToModel(dto)
	case "AzureMetadataService":
		model.AzureMetadata = convertAzureMetadataDTOToModel(dto)
	case "GcpIdentityToken":
		model.GcpIdentity = convertGcpIdentityDTOToModel(dto)
	case "GitHubIdentityToken":
		model.GitHubAction = convertGitHubActionDTOToModel(dto)
	case "Kerberos":
		model.Kerberos = convertKerberosDTOToModel(dto)
	case "KubernetesServiceAccount":
		model.KubernetesService = convertKubernetesDTOToModel(dto)
	case "TerraformIdentityToken":
		model.TerraformWorkspace = convertTerraformDTOToModel(dto)
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

func convertAwsEcsRoleDTOToModel(dto aembit.TrustProviderDTO) *trustProviderAwsEcsRoleModel {
	model := &trustProviderAwsEcsRoleModel{
		AccountID:   types.StringNull(),
		AssumedRole: types.StringNull(),
		RoleARN:     types.StringNull(),
		Username:    types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AwsAccountId":
			model.AccountID = types.StringValue(rule.Value)
		case "AwsAssumedRole":
			model.AssumedRole = types.StringValue(rule.Value)
		case "AwsRoleARN":
			model.RoleARN = types.StringValue(rule.Value)
		case "AwsUsername":
			model.Username = types.StringValue(rule.Value)
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
		Principal: types.StringNull(),
		Realm:     types.StringNull(),
		SourceIP:  types.StringNull(),
	}
	model.AgentControllerIDs = make([]types.String, len(dto.AgentControllerIDs))
	for i, controllerID := range dto.AgentControllerIDs {
		model.AgentControllerIDs[i] = types.StringValue(controllerID)
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

func convertKubernetesDTOToModel(dto aembit.TrustProviderDTO) *trustProviderKubernetesModel {
	decodedKey, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &trustProviderKubernetesModel{
		Issuer:             types.StringNull(),
		Namespace:          types.StringNull(),
		PodName:            types.StringNull(),
		ServiceAccountName: types.StringNull(),
		Subject:            types.StringNull(),
		PublicKey:          types.StringNull(),
		OIDCEndpoint:       types.StringNull(),
	}
	if len(dto.Certificate) > 0 {
		model.PublicKey = types.StringValue(string(decodedKey))
	} else {
		model.OIDCEndpoint = types.StringValue(dto.OidcUrl)
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "KubernetesIss":
			model.Issuer = types.StringValue(rule.Value)
		case "KubernetesIoNamespace":
			model.Namespace = types.StringValue(rule.Value)
		case "KubernetesIoPodName":
			model.PodName = types.StringValue(rule.Value)
		case "KubernetesIoServiceAccountName":
			model.ServiceAccountName = types.StringValue(rule.Value)
		case "KubernetesSub":
			model.Subject = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertGcpIdentityDTOToModel(dto aembit.TrustProviderDTO) *trustProviderGcpIdentityModel {
	model := &trustProviderGcpIdentityModel{
		EMail: types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "Email":
			model.EMail = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertGitHubActionDTOToModel(dto aembit.TrustProviderDTO) *trustProviderGitHubActionModel {
	model := &trustProviderGitHubActionModel{
		Actor:      types.StringNull(),
		Repository: types.StringNull(),
		Workflow:   types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "GithubActor":
			model.Actor = types.StringValue(rule.Value)
		case "GithubRepository":
			model.Repository = types.StringValue(rule.Value)
		case "GithubWorkflow":
			model.Workflow = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertTerraformDTOToModel(dto aembit.TrustProviderDTO) *trustProviderTerraformModel {
	model := &trustProviderTerraformModel{
		OrganizationID: types.StringNull(),
		ProjectID:      types.StringNull(),
		WorkspaceID:    types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "TerraformOrganizationId":
			model.OrganizationID = types.StringValue(rule.Value)
		case "TerraformProjectId":
			model.ProjectID = types.StringValue(rule.Value)
		case "TerraformWorkspaceId":
			model.WorkspaceID = types.StringValue(rule.Value)
		}
	}
	return model
}
