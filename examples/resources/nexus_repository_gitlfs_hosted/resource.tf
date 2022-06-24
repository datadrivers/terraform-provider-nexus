resource "nexus_repository_gitlfs_hosted" "internal" {
  name   = "gitlfs-internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }
}
