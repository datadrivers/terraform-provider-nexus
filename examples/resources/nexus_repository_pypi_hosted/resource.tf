resource "nexus_repository_pypi_hosted" "internal" {
  name   = "pypi-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
