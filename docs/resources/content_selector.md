---
layout: "nexus"
page_title: "Nexus: nexus_content_selector"
subcategory: "Other"
sidebar_current: "docs-nexus-resource-content_selector"
description: |-
  Use this resource to create a Nexus Content Selector
---

# nexus_content_selector

Use this resource to create a Nexus Content Selector

## Example Usage

```hcl
resource "nexus_content_selector" "docker-public" {
  name        = "docker-public"
  description = "A content selector matching public docker images."
  expression  = "path =^ \"/v2/public/\""
}
```

## Argument Reference

The following arguments are supported:

* `expression` - (Required) The content selector expression
* `name` - (Required, ForceNew) Content selector name
* `description` - (Optional) A description of the content selector


