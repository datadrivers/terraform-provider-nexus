---
layout: "nexus"
page_title: "Nexus: nexus_blobstore"
sidebar_current: "docs-nexus-resource-blobstore"
description: |-
  Use this resource to create a Nexus blobstore.
---

# nexus_blobstore

Use this resource to create a Nexus blobstore.

## Example Usage

Example Usage for file store

```hcl
resource "nexus_blobstore" "default" {
  name = "blobstore-file"
  type = "File"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024
    type  = "spaceRemainingQuota"
  }
}
```

Example Usage with S3 bucket

```hcl
resource "nexus_blobstore" "aws" {
  name = "blobstore-s3"
  type = "S3"

  bucket_configuration {
    bucket {
      name   = "aws-bucket-name"
      region = "us-central-1"
    }

    bucket_security {
      access_key_id = "<your-aws-access-key-id>"
      secret_access_key = "<your-aws-secret-access-key>"
    }
  }

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


