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
