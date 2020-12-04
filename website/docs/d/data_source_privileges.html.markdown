---
layout: "nexus"
page_title: "Nexus: data_source_privileges"
sidebar_current: "docs-nexus-data-source"
description: |-
  Retrieve data about Nexus privileges.
---

# data_source_privileges

Retrieve data about Nexus privileges.

## Example Usage

```hcl
data "nexus_privileges" "repository-admin" {
	type = "repository-admin"
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Optional) The privilege type to be queried for.

## Attributes Reference

* `type` - The privilege type
* `name` - The privilege name
* `description` - The description of the privilege
* `readOnly` - Boolean to tell if privilege is read only
