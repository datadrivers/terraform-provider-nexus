---
layout: "nexus"
page_title: "Nexus: nexus_user"
sidebar_current: "docs-nexus-datasource-user"
description: |-
  Use this data source to get a repository data structure
---

# nexus_user

Use this data source to get a repository data structure

## Example Usage

```hcl
data "nexus_user" "admin" {
  userid = "admin"
}
```

## Argument Reference

The following arguments are supported:

* `userid` - (Required) The userid which is required for login

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `email` - The email address associated with the user.
* `firstname` - The first name of the user.
* `lastname` - The last name of the user.
* `roles` - The roles which the user has been assigned within Nexus.
* `status` - The user's status, e.g. active or disabled.


