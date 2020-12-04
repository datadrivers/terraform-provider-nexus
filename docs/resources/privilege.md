---
layout: "nexus"
page_title: "Nexus: nexus_privilege"
sidebar_current: "docs-nexus-resource-privilege"
description: |-
  Use this resource to create a Nexus privilege
---

# nexus_privilege

Use this resource to create a Nexus privilege

## Example Usage

```hcl
resource "nexus_privilege" "example" {
  name    = "example-privilige"
  actions = "all"
  type    = "repository-admin"
}
```

## Argument Reference

The following arguments are supported:

* `actions` - (Required) Actions for the privilege (browse, read, edit, add, delete, all and run)
* `name` - (Required, ForceNew) The name of the privilege
* `type` - (Required) The type of the privilege
* `content_selector` - (Optional) The content selector for the privilege
* `description` - (Optional) A description of the privilege
* `domain` - (Optional) The domain of the privilege
* `format` - (Optional) The format of the privilege
* `pattern` - (Optional) The wildcard privilege pattern
* `repository` - (Optional) The repository of the privilege
* `script_name` - (Optional) The script name related to the privilege


