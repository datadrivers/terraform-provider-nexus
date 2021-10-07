---
layout: "nexus"
page_title: "Nexus: nexus_security_ldap"
subcategory: "Security"
sidebar_current: "docs-nexus-datasource-security_ldap"
description: |-
  Use this data source to work with LDAP security
---

# nexus_security_ldap

Use this data source to work with LDAP security

## Example Usage

```hcl
data "nexus_security_ldap" "default" {}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ldap` - List of ldap configrations
  * `auth_password` - The password to bind with
  * `auth_realm` - The SASL realm to bind to
  * `auth_schema` - Authentication scheme used for connecting to LDAP server
  * `auth_username` - This must be a fully qualified username if simple authentication is used
  * `connection_retry_delay_seconds` - How long to wait before retrying
  * `connection_timeout_seconds` - How long to wait before timeout
  * `group_base_dn` - The relative DN where group objects are found (e.g. ou=Group)
  * `group_id_attribute` - This field specifies the attribute of the Object class that defines the Group ID
  * `group_member_attribute` - LDAP attribute containing the usernames for the group
  * `group_member_format` - The format of user ID stored in the group member attribute
  * `group_object_class` - LDAP class for group objects
  * `group_subtree` - Are groups located in structures below the group base DN
  * `group_type` - Defines a type of groups used: static (a group contains a list of users) or dynamic (a user contains a list of groups)
  * `host` - LDAP server connection hostname
  * `id` - The if of the ldap configuration
  * `ldap_groups_as_roles` - Denotes whether LDAP assigned roles are used as Nexus Repository Manager roles
  * `max_incident_count` - How many retry attempts
  * `name` - LDAP server name
  * `port` - LDAP server connection port to use
  * `protocol` - LDAP server connection Protocol to use
  * `search_base` - LDAP location to be added to the connection URL
  * `use_trust_store` - Whether to use certificates stored in Nexus Repository Manager's truststore
  * `user_base_dn` - The relative DN where user objects are found (e.g. ou=people). This value will have the Search base DN value appended to form the full User search base DN.
  * `user_email_address_attribute` - This is used to find an email address given the user ID
  * `user_id_attribute` - This is used to find a user given its user ID
  * `user_ldap_filter` - LDAP search filter to limit user search
  * `user_member_of_attribute` - Set this to the attribute used to store the attribute which holds groups DN in the user object
  * `user_object_class` - LDAP class for user objects
  * `user_password_attribute` - If this field is blank the user will be authenticated against a bind with the LDAP server
  * `user_real_name_attribute` - This is used to find a real name given the user ID
  * `user_subtree` - Are users located in structures below the user base DN?


