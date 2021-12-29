---
layout: "nexus"
page_title: "Nexus: nexus_repository_yum_hosted"
subcategory: "Repository"
sidebar_current: "docs-nexus-datasource-repository_yum_hosted"
description: |-
  Use this data source to get an existing yum repository
---

# nexus_repository_yum_hosted

Use this data source to get an existing yum repository

## Example Usage

```hcl
data "nexus_repository_yum_hosted" "yummy" {
  name = "yummy"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique identifier for this repository

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cleanup` - Cleanup policies
* `deploy_policy` - Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`
* `online` - Whether this repository accepts incoming requests
* `repodata_depth` - Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5
* `storage` - The storage configuration of the repository
  * `blob_store_name` - Blob store used to store repository contents
  * `strict_content_type_validation` - Whether to validate uploaded content's MIME type appropriate for the repository format
  * `write_policy` - Controls if deployments of and updates to assets are allowed


