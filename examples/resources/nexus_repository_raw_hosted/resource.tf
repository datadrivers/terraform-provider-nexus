resource "nexus_repository_raw_hosted" "internal" {
  name   = "raw-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
