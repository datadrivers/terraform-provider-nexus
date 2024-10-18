# arbitrary example
resource "nexus_blobstore_file" "default" {
  name = "blobstore-file"
  type = "File"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}