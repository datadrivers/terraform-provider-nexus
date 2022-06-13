resource "nexus_repository_nuget_hosted" "internal" {
  name   = "nuget-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
