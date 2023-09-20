resource "nexus_privilege_repository_admin" "example_repository_view" {
	name = "new-repository-admin"
	description = "description"
	actions = ["ADD", "READ", "DELETE", "BROWSE", "EDIT"]
	repository = resource.nexus_repository_helm_hosted.example.name
	format = "helm"
}