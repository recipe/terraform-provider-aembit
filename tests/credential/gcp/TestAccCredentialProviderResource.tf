provider "aembit" {
}

resource "aembit_credential_provider" "gcp" {
	name = "TF Acceptance GCP Workload"
	is_active = true
	google_workload_identity = {
		audience = "audience"
		service_account = "test@test.com"
		lifetime = 1800
	}
}
