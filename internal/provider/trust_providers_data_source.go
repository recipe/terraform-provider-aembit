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
	_ datasource.DataSource              = &trustProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &trustProvidersDataSource{}
)

// NewTrustProvidersDataSource is a helper function to simplify the provider implementation.
func NewTrustProvidersDataSource() datasource.DataSource {
	return &trustProvidersDataSource{}
}

// trustProvidersDataSource is the data source implementation.
type trustProvidersDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *trustProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *trustProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_providers"
}

// Schema defines the schema for the resource.
func (d *trustProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an trust provider.",
		Attributes: map[string]schema.Attribute{
			"trust_providers": schema.ListNestedAttribute{
				Description: "List of trust providers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the trust provider.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the trust provider.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the trust provider.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the trust provider.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"azure_metadata": schema.SingleNestedAttribute{
							Description: "Azure Metadata type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"sku":   schema.StringAttribute{Computed: true},
								"vm_id": schema.StringAttribute{Computed: true},
								"subscription_id": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"aws_ecs_role": schema.SingleNestedAttribute{
							Description: "AWS ECS Role type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Description: "The ID of the AWS account that is hosting the ECS Task.",
									Computed:    true,
								},
								"assumed_role": schema.StringAttribute{
									Description: "The Name of the AWS IAM Role which is running the ECS Task.",
									Computed:    true,
								},
								"role_arn": schema.StringAttribute{
									Description: "The ARN of the AWS IAM Role which is running the ECS Task.",
									Computed:    true,
								},
								"username": schema.StringAttribute{
									Description: "The UsernID of the AWS IAM Account which is running the ECS Task (not commonly used).",
									Computed:    true,
								},
							},
						},
						"aws_metadata": schema.SingleNestedAttribute{
							Description: "AWS Metadata type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"certificate": schema.StringAttribute{
									Computed:    true,
									Description: "PEM Certificate to be used for Signature verification",
								},
								"account_id":                schema.StringAttribute{Computed: true},
								"architecture":              schema.StringAttribute{Computed: true},
								"availability_zone":         schema.StringAttribute{Computed: true},
								"billing_products":          schema.StringAttribute{Computed: true},
								"image_id":                  schema.StringAttribute{Computed: true},
								"instance_id":               schema.StringAttribute{Computed: true},
								"instance_type":             schema.StringAttribute{Computed: true},
								"kernel_id":                 schema.StringAttribute{Computed: true},
								"marketplace_product_codes": schema.StringAttribute{Computed: true},
								"pending_time":              schema.StringAttribute{Computed: true},
								"private_ip":                schema.StringAttribute{Computed: true},
								"ramdisk_id":                schema.StringAttribute{Computed: true},
								"region":                    schema.StringAttribute{Computed: true},
								"version":                   schema.StringAttribute{Computed: true},
							},
						},
						"gcp_identity": schema.SingleNestedAttribute{
							Description: "GCP Identity type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"email": schema.StringAttribute{
									Description: "The Email of the GCP Service Account used by the associated GCP resource.",
									Computed:    true,
								},
							},
						},
						"github_action": schema.SingleNestedAttribute{
							Description: "GitHub Action type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"actor": schema.StringAttribute{
									Description: "The GitHub Actor which initiated the GitHub Action.",
									Computed:    true,
								},
								"repository": schema.StringAttribute{
									Description: "The GitHub Repository associated with the GitHub Action ID Token.",
									Computed:    true,
								},
								"workflow": schema.StringAttribute{
									Description: "The GitHub Workflow execution associated with the GitHub Action ID Token.",
									Computed:    true,
								},
							},
						},
						"kerberos": schema.SingleNestedAttribute{
							Description: "Kerberos type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"agent_controller_ids": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"principal": schema.StringAttribute{Computed: true},
								"realm":     schema.StringAttribute{Computed: true},
								"source_ip": schema.StringAttribute{Computed: true},
							},
						},
						"kubernetes_service_account": schema.SingleNestedAttribute{
							Description: "Kubernetes Service Account type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"issuer": schema.StringAttribute{
									Description: "The Issuer (`iss` claim) of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"namespace": schema.StringAttribute{
									Description: "The Namespace of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"pod_name": schema.StringAttribute{
									Description: "The Pod Name of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"service_account_name": schema.StringAttribute{
									Description: "The Service Account Name of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"subject": schema.StringAttribute{
									Description: "The Subject (`sub` claim) of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"oidc_endpoint": schema.StringAttribute{
									Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"public_key": schema.StringAttribute{
									Description: "The Public Key that can be used to verify the signature of the Kubernetes Service Account Token.",
									Computed:    true,
								},
							},
						},
						"terraform_workspace": schema.SingleNestedAttribute{
							Description: "Terraform Workspace type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"organization_id": schema.StringAttribute{
									Description: "The Organization ID of the calling Terraform Workspace.",
									Computed:    true,
								},
								"project_id": schema.StringAttribute{
									Description: "The Project ID of the calling Terraform Workspace.",
									Computed:    true,
								},
								"workspace_id": schema.StringAttribute{
									Description: "The Workspace ID of the calling Terraform Workspace.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *trustProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state trustProvidersDataSourceModel

	trustProviders, err := d.client.GetTrustProviders(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Trust Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, trustProvider := range trustProviders {
		trustProviderState := convertTrustProviderDTOToModel(ctx, trustProvider)
		state.TrustProviders = append(state.TrustProviders, trustProviderState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
