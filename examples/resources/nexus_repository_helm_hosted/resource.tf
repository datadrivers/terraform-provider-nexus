resource "nexus_repository_helm_hosted" "internal" {
  name   = "helm-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
