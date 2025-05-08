resource "nexus_repository_cargo_hosted" "releases" {
  name   = "cargo-releases"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }

}
