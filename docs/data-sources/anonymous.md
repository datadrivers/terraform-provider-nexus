---
layout: "nexus"
page_title: "Nexus: nexus_anonymous"
subcategory: "Other"
sidebar_current: "docs-nexus-datasource-anonymous"
description: |-
  Use this get the anonymous configuration of the nexus repository manager.
---

# nexus_anonymous

Use this get the anonymous configuration of the nexus repository manager.

## Example Usage

```hcl
data "nexus_anonymous" "nexus" {
}

output "nexus_anonymous_enabled" {
  description = "Anonymous enabled?"
  value       = data.nexus_anonymous.nexus.enabled
}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` - Activate the anonymous access to the repository manager
* `realm_name` - The name of the used realm
* `user_id` - The user id used by anonymous access


