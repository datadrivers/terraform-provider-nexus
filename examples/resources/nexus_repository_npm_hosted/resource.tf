resource "nexus_repository_npm_hosted" "npm" {
  name = "npm"
}

resource "nexus_repository_npm_hosted" "npm1" {
  name = "npm1"
  online = true

  cleanup {
    policy_names = ["policy"]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
