resource "nexus_blobstore_file" "file" {
  name = "blobstore-file"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}
