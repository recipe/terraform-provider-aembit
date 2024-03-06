provider "aembit" {
}

resource "aembit_credential_provider" "aws" {
	name = "TF Acceptance AWS STS"
	is_active = true
	aws_sts = {
		role_arn = "role_arn"
		lifetime = 1800
	}
}
