# nexus_repository_huggingface_proxy

Manages a Hugging Face proxy repository using the dedicated Nexus REST API.
Requires Nexus Repository 3.77 or newer. Uses the common proxy repository fields:
`name`, `online`, `storage`, `proxy`, `negative_cache`, `http_client`, optional
`routing_rule` and `cleanup`. See `repository_cocoapods_proxy` for common fields.

```hcl
resource "nexus_repository_huggingface_proxy" "models" {
  name   = "huggingface-proxy"
  online = true
  storage {
    blob_store_name                = "models"
    strict_content_type_validation = true
  }
  proxy {
    remote_url       = "https://huggingface.co"
    content_max_age  = -1
    metadata_max_age = 1440
  }
  negative_cache {
    enabled = true
    ttl     = 15
  }
  http_client {
    blocked    = false
    auto_block = false
    connection {
      timeout = 120
      retries = 3
    }
  }
}
```

For gated models, configure `http_client.authentication` with `type =
"bearerToken"` and sensitive `bearer_token`. Restrict client read access to the
repository before caching private or gated content.

Import using `tofu import nexus_repository_huggingface_proxy.models huggingface-proxy`.
Clients use `HF_ENDPOINT` and `HF_HUB_DISABLE_XET=1`.
