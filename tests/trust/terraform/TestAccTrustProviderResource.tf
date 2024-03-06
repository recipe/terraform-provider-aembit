provider "aembit" {
}

resource "aembit_trust_provider" "terraform" {
	name = "TF Acceptance Terraform Workspace"
	is_active = true
	terraform_workspace = {
		organization_id = "organization_id"
		project_id = "project_id"
		workspace_id = "workspace_id"
	}
}