---
layout: "nexus"
page_title: "Nexus: nexus_security_realms"
sidebar_current: "docs-nexus-resource-security_realms"
description: |-
  Use this resource to activate Nexus Security LDAP and order
---

# nexus_security_realms

Use this resource to activate Nexus Security LDAP and order

## Example Usage

```hcl
resource "nexus_security_realms" "example" {
  active = [
	"NexusAuthenticatingRealm",
	"NexusAuthorizingRealm",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `active` - (Required) The realm IDs


