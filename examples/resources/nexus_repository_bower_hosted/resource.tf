resource "nexus_repository_bower_hosted" "internal" {
  name   = "bower-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
