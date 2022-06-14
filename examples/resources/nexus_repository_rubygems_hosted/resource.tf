resource "nexus_repository_rubygems_hosted" "internal" {
  name   = "rubygems-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
