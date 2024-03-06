provider "aembit" {
}

resource "aembit_trust_provider" "aws_ecs" {
	name = "TF Acceptance AWS ECS"
	is_active = true
	aws_ecs_role = {
		account_id = "account_id"
		assumed_role = "assumed_role"
		role_arn = "role_arn"
		username = "username"
	}
}