---
layout: "nexus"
page_title: "Nexus: nexus_routing_rule"
subcategory: "Routing"
sidebar_current: "docs-nexus-datasource-routing_rule"
description: |-
  Use this data source to work with routing rules
---

# nexus_routing_rule

Use this data source to work with routing rules

## Example Usage

```hcl
data "nexus_routing_rule" "stop_leaks" {
  name = "stop-leaks"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the routing rule

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - The description of the routing rule
* `matchers` - Matchers is a list of regular expressions used to identify request paths that are allowed or blocked (depending on above mode)
* `mode` - The mode describe how to hande with mathing requests. Possible values: `BLOCK` or `ALLOW`


