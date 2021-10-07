---
layout: "nexus"
page_title: "Nexus: nexus_script"
subcategory: "Other"
sidebar_current: "docs-nexus-resource-script"
description: |-
  Use this resource to create and execute a custom script.
---

# nexus_script

Use this resource to create and execute a custom script.

## Example Usage

```hcl
resource "nexus_script" "repo_pypi_internal" {
  name    = "create-repo-pypi-internal"
  type    = "groovy"
  content = "repository.createPyPiHosted('pypi-internal')"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required) The content of this script.
* `name` - (Required) The name of the script.
* `type` - (Optional) The type of the script. Default: `groovy`


