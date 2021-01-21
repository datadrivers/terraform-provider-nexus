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
    limit = 1024000000
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
* `type` - (Required, ForceNew) The type of the blobstore. Possible values: `S3` or `File`
* `bucket_configuration` - (Optional) The S3 bucket configuration. Needed for blobstore type 'S3'
* `path` - (Optional) The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory
* `soft_quota` - (Optional) Soft quota of the blobstore

The `bucket_configuration` object supports the following:

* `bucket` - (Required) The S3 bucket configuration
* `advanced_bucket_connection` - (Optional) Additional connection configurations
* `bucket_security` - (Optional) Additional security configurations
* `encryption` - (Optional) Additional bucket encryption configurations

The `advanced_bucket_connection` object supports the following:

* `endpoint` - (Optional) A custom endpoint URL for third party object stores using the S3 API.
* `force_path_style` - (Optional) Setting this flag will result in path-style access being used for all requests.
* `signer_type` - (Optional) An API signature version which may be required for third party object stores using the S3 API.

The `bucket` object supports the following:

* `name` - (Required) The name of the S3 bucket
* `region` - (Required) The AWS region to create a new S3 bucket in or an existing S3 bucket's region
* `expiration` - (Optional) How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)
* `prefix` - (Optional) The S3 blob store (i.e S3 object) key prefix

The `bucket_security` object supports the following:

* `access_key_id` - (Required) An IAM access key ID for granting access to the S3 bucket
* `secret_access_key` - (Required) The secret access key associated with the specified IAM access key ID
* `role` - (Optional) An IAM role to assume in order to access the S3 bucket
* `session_token` - (Optional) An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket

The `encryption` object supports the following:

* `encryption_key` - (Optional) The encryption key.
* `encryption_type` - (Optional) The type of S3 server side encryption to use.

The `soft_quota` object supports the following:

* `limit` - (Required) The limit in Bytes. Minimum value is 1000000
* `type` - (Required) The type to use such as spaceRemainingQuota, or spaceUsedQuota

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_space_in_bytes` - Available space in Bytes
* `blob_count` - Count of blobs
* `total_size_in_bytes` - The total size of the blobstore in Bytes


