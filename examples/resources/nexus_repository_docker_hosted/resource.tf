resource "nexus_repository_docker_hosted" "example" {
  name   = "example"
  online = true

  docker {
    force_basic_auth = false
    v1_enabled       = false
    subdomain        = "docker"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
