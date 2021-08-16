---
layout: "nexus"
page_title: "Nexus: nexus_anonymous"
sidebar_current: "docs-nexus-resource-anonymous"
description: |-
  Use this resource to change the anonymous configuration of the nexus repository manager.
---

# nexus_anonymous

Use this resource to change the anonymous configuration of the nexus repository manager.

## Example Usage

```hcl
resource "nexus_anonymous" "example" {
  enabled = true
  user_id = "exampleUser"
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Optional) Activate the anonymous access to the repository manager. Default: false
* `realm_name` - (Optional) The name of the used realm. Default: "NexusAuthorizingRealm"
* `user_id` - (Optional) The user id used by anonymous access. Default: "anonymous"


