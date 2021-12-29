---
layout: "nexus"
page_title: "Nexus: nexus_repository_yum_hosted"
subcategory: "Repository"
sidebar_current: "docs-nexus-datasource-repository_yum_hosted"
description: |-
  Use this data source to get a repository data structure
---

# nexus_repository_yum_hosted

Use this data source to get a repository data structure

## Example Usage

```hcl
data "nexus_repository_yum_hosted" "yummy" {
  name = "yummy"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique identifier for this repository
* `cleanup` - (Optional) Cleanup policies
* `deploy_policy` - (Optional) Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`
* `online` - (Optional) Whether this repository accepts incoming requests
* `repodata_depth` - (Optional) Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5
* `storage` - (Optional) The storage configuration of the repository

The `cleanup` object supports the following:

* `policy_names` - (Optional) List of policy names

The `storage` object supports the following:

* `blob_store_name` - (Optional) Blob store used to store repository contents
* `strict_content_type_validation` - (Optional) Whether to validate uploaded content's MIME type appropriate for the repository format
* `write_policy` - (Optional) Controls if deployments of and updates to assets are allowed

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - Repository type


