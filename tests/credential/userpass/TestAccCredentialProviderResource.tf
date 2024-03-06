provider "aembit" {
}

resource "aembit_credential_provider" "userpass" {
	name = "TF Acceptance Username Password"
	is_active = true
	username_password = {
		username = "username"
		password = "password"
	}
}
