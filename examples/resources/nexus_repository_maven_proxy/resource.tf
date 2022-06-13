resource "nexus_repository_maven_proxy" "maven_central" {
  name   = "maven-central-repo1"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = "https://repo1.maven.org/maven2/"
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

  maven {
    version_policy = "RELEASE"
    layout_policy  = "PERMISSIVE"
  }
}
