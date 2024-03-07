# Aembit Terraform Provider

This is the repository for the Aembit Cloud Terraform Provider. Learn more about Aembit at https://aembit.io/

## Support, Bugs, Feature Requests

Any requests should be filed under the Issues section of this repository. All filed issues will be handled on a "best effort" basis.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.6
- [Go](https://golang.org/doc/install) >= 1.20

## Getting Started

The provider can be installed by running `terraform init`.

The provider block can be specified as follows:
```shell
terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}


provider "aembit" {
}
```