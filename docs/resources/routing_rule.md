---
layout: "nexus"
page_title: "Nexus: nexus_routing_rule"
sidebar_current: "docs-nexus-resource-routing_rule"
description: |-
  Use this resource to create a Nexus Routing Rule.
---

# nexus_routing_rule

Use this resource to create a Nexus Routing Rule.

## Example Usage

```hcl
resource "nexus_routing_rule" "stop_leaks" {
  name        = "stop-leaks"
  description = "Prevent requests of internal names"
  mode        = "BLOCK"
  matchers    = [
	"^/com/example/.*",
	"^/org/example/.*",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `matchers` - (Required) Matchers is a list of regular expressions used to identify request paths that are allowed or blocked (depending on above mode)
* `name` - (Required, ForceNew) The name of the routing rule
* `description` - (Optional) The description of the routing rule
* `mode` - (Optional) The mode describe how to hande with mathing requests. Possible values: `BLOCK` or `ALLOW` Default: `BLOCK`


