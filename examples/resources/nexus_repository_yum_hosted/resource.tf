resource "nexus_repository_yum_hosted" "yum" {
  name = "yummy"
}

resource "nexus_repository_yum_hosted" "yum1" {
  deploy_policy  = "STRICT"
  name = "yummy1"
  online = true
  repodata_depth = 4

  cleanup {
    policy_names = ["policy"]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
