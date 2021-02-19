---
layout: "nexus"
page_title: "Nexus: nexus_security_user"
sidebar_current: "docs-nexus-resource-security_user"
description: |-
  Use this resource to manage users
---

# nexus_security_user

Use this resource to manage users

## Example Usage

```hcl
resource "nexus_security_user" "admin" {
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


