---
layout: "nexus"
page_title: "Nexus: nexus_security_ldap_order"
sidebar_current: "docs-nexus-resource-security_ldap_order"
description: |-
  Use this resource to change the LDAP order.
---

# nexus_security_ldap_order

Use this resource to change the LDAP order.

## Example Usage

```hcl
resource "nexus_security_ldap_order" "order" {
	order = ["127.0.0.1", "localhost"]
}

```

## Argument Reference

The following arguments are supported:




