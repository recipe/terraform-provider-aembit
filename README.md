# Aembit Terraform Provider

This repository houses the Aembit Terraform Provider source. This code is under development.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Getting Started

Since both the Aembit API client code and the Terraform provider code are under development, we
will need to use local copies of these repositories for development.

### Clone the Aembit API client code

```shell
cd ~/src
git clone git@github.com:Aembit/aembit_api_client.git
```

### Override the Golang module dependency

Since our API Client code is not published, we will use a local copy of the library, and tell
our Terraform code to use a local copy of the module.

Assuming that both the Terraform provider repository and the API client library repository are
both cloned under `~src`, redirect the Aembit API client module reference as follows:

```shell
go mod edit -replace aembit.io/aembit=../aembit_api_client
```

This command will modify `go.mod`; you will see a line like the following:

> replace aembit.io/aembit => ../aembit_api_client

After modifying go.mod, retrieve the local Aembit API client module.

```shell
go get aembit.io/aembit
```

You should see output like the following:

> go: added aembit.io/aembit v0.0.0-00010101000000-000000000000

### Modify Terraform configuration to use local Provider code

Since our Terraform provider code is not published, we will tell our Terraform binary to use the
local provider code.

Edit `~/.terraformrc` so that it looks like the following:

```
provider_installation {

  dev_overrides {
      "aembit.io/dev/aembit" = "/Users/jkwon/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Replace the Golang path with the path appropriate for your system.

## Building The Provider

Most provider-related files will be under the `internal/provider` directory. After
adding new files to that directory, or modifying files under that directory, build
the provider.

1. `cd` to the Terraform provider repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

Besides the Aembit API client library (which is local), if you need to use any
additional modules:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`. However, do <b>not</b> commit
the override changes for the `aembit_api_client` library.

## Testing the provider

Example Terraform templates are under the `examples` directory. For example:

```shell
cd examples/server-workloads
terraform plan
```

To test additional resources and datasources, create new directories under the `examples`
and implement your Terraform template there.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
