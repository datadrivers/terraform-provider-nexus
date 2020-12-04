---
layout: "nexus"
page_title: "Nexus: nexus_role"
sidebar_current: "docs-nexus-resource-role"
description: |-
  Use this resource to create a Nexus Role.
---

# nexus_role

Use this resource to create a Nexus Role.

## Example Usage

```hcl
resource "nexus_role" "docker_deploy" {
  description = "Docker deployment role"
  name        = "docker-deploy"
  privileges = [
    "nx-repository-view-docker-*-*",
  ]
  roleid = "docker-deploy"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.
* `roleid` - (Required, ForceNew) The id of the role.
* `description` - (Optional) The description of this role.
* `privileges` - (Optional) The privileges of this role.
* `roles` - (Optional) The roles of this role.


