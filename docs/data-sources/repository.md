---
layout: "nexus"
page_title: "Nexus: nexus_repository"
sidebar_current: "docs-nexus-datasource-repository"
description: |-
  Use this data source to get a repository data structure
---

# nexus_repository

Use this data source to get a repository data structure

## Example Usage

```hcl
data "nexus_repository" "maven-central" {
  name = "maven-central"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique identifier for this repository


