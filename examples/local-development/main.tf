resource "nexus_repository_npm_proxy" "repo" {
  name   = "private_repo"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }

  proxy {
    remote_url       = "https://gitlab.com/api/v4/projects/70769722/packages/npm/"
    content_max_age  = 1338
    metadata_max_age = 1440
  }

  negative_cache {
    enabled = false
    ttl     = 1440
  }

  http_client {
    blocked    = false
    auto_block = true
    authentication {
      type         = "bearerToken"
      bearer_token = "test-token"
    }
  }



  # # we need to ignore the authentication block,
  # # because it is not supported by the provider because of nexus API:
  # # https://github.com/datadrivers/terraform-provider-nexus/issues/158
  # # https://github.com/sonatype/nexus-public/issues/247
  # lifecycle {
  #   ignore_changes = [
  #     http_client.0.authentication,
  #   ]
  # }
}
