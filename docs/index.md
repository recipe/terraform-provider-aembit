---
layout: ""
page_title: "Provider: Aembit Cloud"
description: |-
  This Aembit Cloud provider provides resources and data sources to manage the Aembit platform as infrastructure-as-code, through the Aembit management API.
---

# Aembit Provider

The Aembit provider interacts with the configuration of the Aembit platform via the management API. The provider requires credentials before it can be used.

## Getting Started

To get started using the Aembit Terraform provider, first you'll need an active Aembit cloud tenant.  Get instant access with a [Aembit trial account](https://useast2.aembit.io/signup), or read more about Aembit at [aembit.io](https://aembit.io)

## Provider Authentication

### Authenticate using Aembit native authentication

Aembit supports authentication to the Aembit API using a native authentication capability which utilizes OIDC (Open ID Connect tokens) ID Tokens. This capability requires configuring your Aembit tenant with the appropriate components as follows:
* **Client Workload:** This workload identifies the execution environment of the Terraform Provider, either in Terraform Cloud, GitHub Actions, or another Aembit-supported Serverless platform.
* **Trust Provider:** This component ensures the authentication of the Client Workload using attestation of the platform ID Token and associated match rules.
  * Match Rules can be configured for platform-specific restrictions, for example repository on GitHub or workspace ID on Terraform Cloud.
* **Credential Provider:** This associates the Client Workload with an Aembit Role to ensure that the Client Workload has access to only the applicable Aembit resources.
  * Note: The Aembit API hostname will be provided as an Audience value here and can be copied to the Server Workload hostname field.
* **Server Workload:** This workload identifies the Aembit tenant-specific API endpoint.
* **Access Policy:** This policy associates the previously configured components and ensures that only this specific workload has the intended access as defined.

After configuring these Aembit resources, the Client ID from the Trust Provider can be configured for the Aembit Terraform Provider, enabling automatic native authentication for the configured Workload.
The Client ID can be configured using the `client_id` field in the Aembit provider configuration block or with the `AEMBIT_CLIENT_ID` environment variable.

<div style="background: #d1ecf1; padding: 0.75rem 1.25rem; margin: 0 0 1rem 0; border-radius: 8px;">:grey_exclamation: <b>Terraform Cloud Configuration</b>
<br>Setting the environment variable TFC_WORKLOAD_IDENTITY_AUDIENCE is required for Terraform Cloud Workspace ID Tokens. The value for this variable will be provided by your Aembit Cloud tenant Trust Provider and references your tenant-specific endpoint.</div>

#### Sample Terraform Config

```terraform
terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
  # This client_id configuration may be set here or in the AEMBIT_CLIENT_ID environment variable.
  # Note: This is a sample value and must be replaced with the Aembit Trust Provider generated value.
  client_id = "aembit:useast2:tenant:identity:github_idtoken:0bc4dbcd-e9c8-445b-ac90-28f47b8649cc"
}

resource "aembit_client_workload" "client" {
  # Resource configuration
}
```

```shell
$ terraform plan
```

### Authenticate using an environment variable access token

```terraform
terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}


provider "aembit" {
}

resource "aembit_client_workload" "client" {
  # Resource configuration
}
```

```shell
$ export AEMBIT_TENANT_ID="tenant"
$ export AEMBIT_TOKEN="token-from-console"
$ terraform plan
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `client_id` (String) The Aembit Trust Provider Client ID to use for authentication to the Aembit Cloud Tenant instance (recommended).
- `tenant` (String) Tenant ID of the specific Aembit Cloud instance.
- `token` (String, Sensitive) Access Token to use for authentication to the Aembit Cloud Tenant instance.

