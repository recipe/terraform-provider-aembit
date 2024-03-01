provider "aembit" {
}

resource "aembit_integration" "wiz" {
	name = "TF Acceptance Wiz"
	is_active = true
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://endpoint"
	oauth_client_credentials = {
		token_url = "https://url/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

resource "aembit_access_condition" "wiz" {
	name = "TF Acceptance Wiz"
	is_active = true
	integration_id = aembit_integration.wiz.id
	wiz_conditions = {
		max_last_seen = 3600
		container_cluster_connected = true
	}
}