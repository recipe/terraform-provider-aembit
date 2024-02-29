provider "aembit" {
}

resource "aembit_credential_provider" "api_key" {
	name = "TF Acceptance API Key"
	api_key = {
		api_key = "test_api_key"
	}
}
