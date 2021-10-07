---
layout: "nexus"
page_title: "Nexus: nexus_security_ldap_order"
subcategory: "Security"
sidebar_current: "docs-nexus-resource-security_ldap_order"
description: |-
  Use this resource to change the LDAP order.
---

# nexus_security_ldap_order

Use this resource to change the LDAP order.

## Example Usage

Set order of LDAP server

```hcl
resource "nexus_security_ldap" "server1" {
  ...
  name = "server1"
}

resource "nexus_security_ldap" "server2" {
  ...
  name = "server2"
}

resource "nexus_security_ldap_order" {
  order = [
    nexus_security_ldap.server1.name,
    nexus_security_ldap.server2.name,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `order` - (Required) Ordered list of LDAP server


