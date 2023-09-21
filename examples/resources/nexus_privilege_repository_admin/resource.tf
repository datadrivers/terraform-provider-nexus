resource "nexus_repository_helm_hosted" "example" {
  name   = "example_repository"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}

resource "nexus_privilege_repository_admin" "example" {
  name        = "example_privilege"
  description = "description"
  actions     = ["ADD", "READ", "DELETE", "BROWSE", "EDIT"]
  repository  = resource.nexus_repository_helm_hosted.example.name
  format      = "helm"
}
