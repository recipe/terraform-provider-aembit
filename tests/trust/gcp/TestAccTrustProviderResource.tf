provider "aembit" {
}

resource "aembit_trust_provider" "gcp" {
	name = "TF Acceptance GCP Identity"
	is_active = true
	gcp_identity = {
		email = "email"
	}
}