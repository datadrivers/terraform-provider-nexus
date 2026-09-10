# arbitrary example
resource "nexus_blobstore_file" "default" {
  name = "blobstore-file"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}

resource "nexus_http_client" "system_proxy" {
  timeout    = 30
  retries    = 3
  user_agent = "Nexus Repository Manager/Custom"

  non_proxy_hosts = [
    "localhost",
    "127.0.0.1",
    "*.internal.lan"
  ]

  http_proxy {
    enabled = true
    host    = "squid.internal.lan"
    port    = 3128
  }

  https_proxy {
    enabled = false
  }
}
