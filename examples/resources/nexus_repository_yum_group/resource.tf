resource "nexus_repository_yum_hosted" "internal" {
  name = "internal"

  deploy_policy  = "STRICT"
  repodata_depth = 4

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}

resource "nexus_repository_yum_group" "group" {
  name   = "yum-group"
  online = true

  group {
    member_names = [
      nexus_repository_yum_hosted.internal.name,
    ]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
