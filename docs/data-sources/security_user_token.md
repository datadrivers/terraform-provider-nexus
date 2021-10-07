---
layout: "nexus"
page_title: "Nexus: nexus_security_user_token"
subcategory: "Security"
sidebar_current: "docs-nexus-datasource-security_user_token"
description: |-
  Use this data source to get the global user-token configuration.
---

# nexus_security_user_token

Use this data source to get the global user-token configuration.

---
**PRO Feature**
---

## Example Usage

```hcl
data "nexus_security_user_token" "nexus" {}

output "nexus_user_token_enabled" {
  description = "User Tokens enabled?"
  value       = data.nexus_security_user_token.nexus.enabled
}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` - Activate the feature of user tokens.
* `protect_content` - Require user tokens for repository authentication. This does not effect UI access.


