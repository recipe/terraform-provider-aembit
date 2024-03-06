provider "aembit" {
}

resource "aembit_credential_provider" "snowflake" {
	name = "TF Acceptance Snowflake Token"
	is_active = true
	snowflake_jwt = {
		account_id = "account_id"
		username = "username"
	}
}
