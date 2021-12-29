---
layout: "nexus"
page_title: "Nexus: nexus_repository_yum_hosted"
subcategory: "Repository"
sidebar_current: "docs-nexus-resource-repository_yum_hosted"
description: |-
  Use this resource to create a hosted yum repository.
---

# nexus_repository_yum_hosted

Use this resource to create a hosted yum repository.

## Example Usage

```hcl

resource "nexus_repository_yum_hosted" "yum" {
  name = "yummy"
}

resource "nexus_repository_yum_hosted" "yum1" {
  deploy_policy  = "STRICT"
  name = "yummy1"
  online = true
  repodata_depth = 4

  cleanup {
    policy_names = ["policy"]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique identifier for this repository
* `storage` - (Required) The storage configuration of the repository
* `cleanup` - (Optional) Cleanup policies
* `deploy_policy` - (Optional) Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`
* `online` - (Optional) Whether this repository accepts incoming requests
* `repodata_depth` - (Optional) Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5

The `cleanup` object supports the following:

* `policy_names` - (Optional) List of policy names

The `storage` object supports the following:

* `blob_store_name` - (Optional) Blob store used to store repository contents
* `strict_content_type_validation` - (Optional) Whether to validate uploaded content's MIME type appropriate for the repository format
* `write_policy` - (Optional) Controls if deployments of and updates to assets are allowed

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - Repository type


