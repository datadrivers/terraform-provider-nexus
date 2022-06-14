resource "nexus_repository_r_hosted" "internal" {
  name   = "r-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
