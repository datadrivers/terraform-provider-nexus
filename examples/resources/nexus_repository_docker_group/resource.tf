resource "nexus_repository_docker_hosted" "internal" {
  name = "internal"

  docker {
    force_basic_auth = false
    v1_enabled       = false
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}

resource "nexus_repository_docker_proxy" "dockerhub" {
  name = "dockerhub"

  docker {
    force_basic_auth = false
    v1_enabled       = false
  }

  docker_proxy {
    index_hub = "HUB"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = "https://registry-1.docker.io"
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled      = true
    time_to_live = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_docker_group" "group" {
  name   = "docker-group"
  online = true

  docker {
    force_basic_auth = false
    http_port        = 8080
    https_port       = 8433
    v1_enabled       = false
  }

  group {
    member_names = [
      nexus_repository_docker_hosted.internal.name,
      nexus_repository_docker_proxy.dockerhub.name
    ]
    writable_member = nexus_repository_docker_hosted.internal.name
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
