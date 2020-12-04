---
layout: "nexus"
page_title: "Nexus: data_source_user"
sidebar_current: "docs-nexus-data-source"
description: |-
  Retrieve data about a Nexus user.
---

# data_source_user

Retrieve data about a Nexus user.

## Example Usage

```hcl
data "nexus_user" "admin" {
	userid = admin
}
```

## Arguments Reference

- `userid` - The user's user id
## Attributes Reference

- `userId` - The user's user id
- `firstName` - The user's first name
- `lastName` - The user's last name
- `emailAddress` - The user's email address
- `source` - The user's source
- `status` - The user's status
- `readOnly` - A boolean value which indicates whether or not the user is read only
- `roles` - A list containing any user roles
- `externalRoles` - A list containing any external roles
