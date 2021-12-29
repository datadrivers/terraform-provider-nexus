---
layout: "nexus"
page_title: "Nexus: nexus_blobstore_file"
subcategory: "Blobstore"
sidebar_current: "docs-nexus-datasource-blobstore_file"
description: |-
  Use this to read an existing Nexus file blobstore.
---

# nexus_blobstore_file

Use this to read an existing Nexus file blobstore.

## Example Usage

```hcl
data "nexus_blobstore_file" "default" {
	name = "default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Blobstore name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_space_in_bytes` - Available space in Bytes
* `blob_count` - Count of blobs
* `path` - The path to the blobstore contents
* `soft_quota` - Soft quota of the blobstore
  * `limit` - The limit in Bytes. Minimum value is 1000000
  * `type` - The type to use such as spaceRemainingQuota, or spaceUsedQuota
* `total_size_in_bytes` - The total size of the blobstore in Bytes


