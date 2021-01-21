---
layout: "nexus"
page_title: "Nexus: nexus_security_realms"
sidebar_current: "docs-nexus-datasource-security_realms"
description: |-
  Use this data source to work with Realms security
---

# nexus_security_realms

Use this data source to work with Realms security

## Example Usage

```hcl
data "nexus_security_realms" "default" {}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `active` - List of active security realms
  * `id` - Realm ID
  * `name` - Realm name
* `available` - List of available security realms
  * `id` - Realm ID
  * `name` - Realm name


