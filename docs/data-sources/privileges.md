---
layout: "nexus"
page_title: "Nexus: nexus_privileges"
subcategory: "Other"
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

* `domain` - (Optional, ForceNew) The domain of the privilege
* `format` - (Optional, ForceNew) The format of the privilege
* `name` - (Optional, ForceNew) The name of the privilege
* `repository` - (Optional, ForceNew) The repository of the privilege
* `type` - (Optional, ForceNew) The type of the privilege

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `privileges` - List of privileges
  * `actions` - Actions for the privilege (browse, read, edit, add, delete, all and run)
  * `content_selector` - The content selector for the privilege
  * `description` - A description of the privilege
  * `domain` - The domain of the privilege
  * `format` - The format of the privilege
  * `name` - The name of the privilege
  * `pattern` - The wildcard privilege pattern
  * `read_only` - Indicates whether the privilege can be changed. External values supplied to this will be ignored by the system.
  * `repository` - The repository of the privilege
  * `type` - The type of the privilege


