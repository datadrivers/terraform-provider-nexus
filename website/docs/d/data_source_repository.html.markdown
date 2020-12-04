---
layout: "nexus"
page_title: "Nexus: data_source_repository"
sidebar_current: "docs-nexus-data-source"
description: |-
  Retrieve data about an existing repository.
---

# data_source_repository

  Retrieve data about an existing repository.

## Example Usage

```hcl
data "nexus_repository" "maven-public" {
  name = "maven-public"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Repository name.
* 
## Attributes Reference

* `name` - The repository name
* `format` - The format (kind) of the repository
* `type` - The repository's type
* `url` - The URL of the repository
* `attribute` - A map of attributes assigned to the repository
