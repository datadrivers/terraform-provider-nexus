---
layout: "nexus"
page_title: "Nexus: nexus_blobstore"
sidebar_current: "docs-nexus-resource-blobstore"
description: |-
  Use this resource to get create a Nexus blobstore.
---

# nexus_blobstore

Use this resource to get create a Nexus blobstore.

## Example Usage

```hcl
resource "nexus_blobstore" "default" {
  name = "blobstore-01"
  type = "File"
  path = "/nexus-data/blobstore-01"

  soft_quota {
    limit = 1024
    type  = "spaceRemainingQuota"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Blobstore name
* `path` - (Optional) The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory


