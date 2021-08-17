---
layout: "nexus"
page_title: "Nexus: nexus_security_user_token"
sidebar_current: "docs-nexus-resource-security_user_token"
description: |-
  Use this resource to manage the global configuration for the user-tokens
---

# nexus_security_user_token

Use this resource to manage the global configuration for the user-tokens

---
**PRO Feature**
---

## Example Usage

```hcl
resource "nexus_security_user_token" "nexus" {
    enabled         = true
	protect_content = false
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required) Activate the feature of user tokens.
* `protect_content` - (Optional) Require user tokens for repository authentication. This does not effect UI access.


