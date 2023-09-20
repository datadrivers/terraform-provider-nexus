resource "nexus_privilege_repository_admin" "example" {
  name        = "example_privilege"
  description = "description"
  actions     = ["ADD", "READ", "DELETE", "BROWSE", "EDIT"]
  repository  = resource.nexus_repository_helm_hosted.example.name
  format      = "helm"
}
