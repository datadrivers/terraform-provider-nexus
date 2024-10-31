resource "nexus_blobstore_file" "blobby" {

  name = "blobby"
  path = "/nexus-data/blobby"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}

resource "nexus_repository_docker_hosted" "example" {
  name = "example-container"

  storage {
    blob_store_name = "default"
    strict_content_type_validation = false
  }

  docker {
    v1_enabled = false
    force_basic_auth = false
  }
}

# resource "nexus_security_cleanup_policy" "cleanup" {
#   name = "CleanerUpper"
# }