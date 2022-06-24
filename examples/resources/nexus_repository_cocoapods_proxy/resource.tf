resource "nexus_repository_cocoapods_proxy" "cocoapods_org" {
  name   = "cocoapods-org"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = "https://cocoapods.org/"
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
