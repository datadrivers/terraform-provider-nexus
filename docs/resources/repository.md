---
layout: "nexus"
page_title: "Nexus: nexus_repository"
sidebar_current: "docs-nexus-resource-repository"
description: |-
  Use this resource to create a Nexus Repository.
---

# nexus_repository

Use this resource to create a Nexus Repository.

## Example Usage

Example Usage - Apt hosted repository

```hcl
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
```

Example Usage - Docker group repository

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `format` - (Required, ForceNew) Repository format
* `name` - (Required) A unique identifier for this repository
* `type` - (Required, ForceNew) Repository type
* `online` - (Optional) Whether this repository accepts incoming requests


