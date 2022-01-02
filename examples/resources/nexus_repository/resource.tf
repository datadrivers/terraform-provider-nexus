# Example Usage - Apt hosted repository
resource "nexus_repository" "apt_hosted" {
  name   = "apt-repo"
  format = "apt"
  type   = "hosted"

  apt {
    distribution = "bionic"
  }

  apt_signing {
    keypair    = "<keypair>"
    passphrase = "<passphrase>"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}

# Example Usage - Docker group repository
resource "nexus_repository" "docker_group" {
  name   = "docker-group"
  format = "docker"
  type   = "group"
  online = true

  group {
    member_names = [
      "docker_releases",
      "docker_hub"
    ]
  }

  docker {
    force_basic_auth = false
    http_port        = 5000
    https_port       = 0
    v1enabled        = false
  }

  storage {
    blob_store_name                = "docker_group_blobstore"
    strict_content_type_validation = true
  }
}
