---
layout: "nexus"
page_title: "Nexus: data_source_blobstore"
sidebar_current: "docs-nexus-data-source"
description: |-
  Retrieve data about an existing blobstore.
---

# data_source_blobstore

Retrieve information about an existing blobstore.

## Example Usage

```hcl
data "nexus_blobstore" "default" {
  name = "default"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Blobstore name.

## Attributes Reference

* `available_space_in_bytes` - Available size in bytes
* `blob_count` - Amount of blobs inside
* `name` - Blobstore name
* `soft_quota` - Softquotas set for the blobstore.
* `soft_quota.limit` - The limit in MB.
* `soft_quota.type` - The type to use such as spaceRemainingQuota, or spaceUsedQuota.
* `total_size_in_bytes` - Size in bytes
* `type` - Blobstore Type such as `File` or `S3`