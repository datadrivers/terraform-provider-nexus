---
layout: "nexus"
page_title: "Nexus: nexus_content_selector"
sidebar_current: "docs-nexus-resource-content_selector"
description: |-
  Use this resource to create a Nexus Content Selector
---

# nexus_content_selector

Use this resource to create a Nexus Content Selector

## Example Usage

```hcl
resource "nexus_content_selector" "selector" {
	description = "My selector"
	expression  = "format == \"maven2\" and path =^ \"/org/sonatype/nexus\""
	name        = "selector"
}
```

## Argument Reference

The following arguments are supported:

* `expression` - (Required) The content selector expression
* `name` - (Required, ForceNew) Content selector name
* `description` - (Optional) A description of the content selector


