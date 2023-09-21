resource "nexus_repository_helm_hosted" "example" {
  name   = "example_repository"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}

resource "nexus_security_content_selector" "example" {
  name        = "example_content_selector"
  description = "A content selector matching public docker images."
  expression  = "path =^ \"/v2/public/\""
}

resource "nexus_privilege_repository_content_selector" "example" {
  name             = "example_privilege"
  description      = "description"
  actions          = ["ADD", "READ", "DELETE", "BROWSE", "EDIT"]
  repository       = resource.nexus_repository_helm_hosted.example.name
  format           = "helm"
  content_selector = resource.nexus_security_content_selector.example.name
}
