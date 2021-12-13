---
layout: "nexus"
page_title: "Nexus: nexus_user"
subcategory: "Other"
sidebar_current: "docs-nexus-resource-user"
description: |-
  This resource is deprecated. Please use the data source "nexus_security_user" instead.
---

# nexus_user

This resource is deprecated. Please use the data source "nexus_security_user" instead.

Use this resource to manage users

## Example Usage

```hcl
resource "nexus_user" "admin" {
  userid    = "admin"
  firstname = "Administrator"
  lastname  = "User"
  email     = "nexus@example.com"
  password  = "admin123"
  roles     = ["nx-admin"]
  status    = "active"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email address associated with the user.
* `firstname` - (Required) The first name of the user.
* `lastname` - (Required) The last name of the user.
* `password` - (Required) The password for the user.
* `userid` - (Required, ForceNew) The userid which is required for login. This value cannot be changed.
* `roles` - (Optional) The roles which the user has been assigned within Nexus.
* `status` - (Optional) The user's status, e.g. active or disabled.


