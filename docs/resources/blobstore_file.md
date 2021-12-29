---
layout: "nexus"
page_title: "Nexus: nexus_blobstore_file"
subcategory: "Blobstore"
sidebar_current: "docs-nexus-resource-blobstore_file"
description: |-
  Use this resource to create a Nexus file blobstore.
---

# nexus_blobstore_file

Use this resource to create a Nexus file blobstore.

## Example Usage

```hcl
resource "nexus_blobstore_file" "default" {
  name = "blobstore-file"
  type = "File"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Blobstore name
* `path` - (Optional) The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory
* `soft_quota` - (Optional) Soft quota of the blobstore

The `soft_quota` object supports the following:

* `limit` - (Required) The limit in Bytes. Minimum value is 1000000
* `type` - (Required) The type to use such as spaceRemainingQuota, or spaceUsedQuota

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_space_in_bytes` - Available space in Bytes
* `blob_count` - Count of blobs
* `total_size_in_bytes` - The total size of the blobstore in Bytes


