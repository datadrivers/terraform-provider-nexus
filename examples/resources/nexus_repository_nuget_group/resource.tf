resource "nexus_repository_nuget_hosted" "internal" {
  name   = "internal"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}

resource "nexus_repository_nuget_proxy" "nuget_org" {
  name   = "nuget-org"
  online = true

  nuget_version            = "V3"
  query_cache_item_max_age = 3600

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = "https://api.nuget.org/v3/index.json"
    content_max_age  = 1440
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
  }
}

resource "nexus_repository_nuget_group" "group" {
  name   = "nuget-group"
  online = true

  group {
    member_names = [
      nexus_repository_nuget_hosted.internal.name,
      nexus_repository_nuget_proxy.nuget_org.name,
    ]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
