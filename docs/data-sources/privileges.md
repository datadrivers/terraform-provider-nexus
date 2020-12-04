---
layout: "nexus"
page_title: "Nexus: nexus_privileges"
sidebar_current: "docs-nexus-datasource-privileges"
description: |-
  Use this data source to work with privileges
---

# nexus_privileges

Use this data source to work with privileges

## Example Usage

```hcl
data "nexus_privileges" "example" {
  domain     = "application"
  format     = "maven2"
  repository = "maven-public"
  type       = "repository-admin"
}
```

## Argument Reference

The following arguments are supported:




