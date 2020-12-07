---
layout: "nexus"
page_title: "Nexus: nexus_security_ldap"
sidebar_current: "docs-nexus-resource-security_ldap"
description: |-
  Use this resource to create a Nexus Security LDAP
---

# nexus_security_ldap

Use this resource to create a Nexus Security LDAP

## Example Usage

```hcl
resource "nexus_security_ldap" "example" {
  auth_password                  = "t0ps3cr3t"
  auth_realm                     = "EXAMPLE"
  auth_schema                    = ""
  auth_username                  = "admin"
  connection_retry_delay_seconds = 1
  connection_timeout_seconds     = 1
  group_base_dn                  = "ou=Group"
  group_id_attribute             = "cn"
  group_member_attribute         = "memberUid"
  group_member_format            = "uid=${username},ou=people,dc=example,dc=com"
  group_object_class             = "example"
  group_subtree                  = true
  host                           = "ldap.example.com"
  ldap_groups_as_roles           = true
  max_incident_count             = 1
  name                           = "example-ldap"
  port                           = 389
  protocol                       = "LDAP"
  search_base                    = "dc=example,dc=com"
  user_base_dn                   = "ou=people"
  user_email_address_attribute   = "mail"
  user_id_attribute              = "uid"
  user_ldap_filter               = "(|(mail=*@example.com)(uid=dom*))"
  user_member_of_attribute       = "memberOf"
  user_object_class              = "posixGroup"
  user_password_attribute        = "exmaple"
  user_real_name_attribute       = "cn"
  user_subtree                   = true
}
```

## Argument Reference

The following arguments are supported:

* `auth_schema` - (Required) Authentication scheme used for connecting to LDAP server
* `auth_username` - (Required) This must be a fully qualified username if simple authentication is used. Required if authScheme other than none.
* `connection_retry_delay_seconds` - (Required) How long to wait before retrying
* `connection_timeout_seconds` - (Required) How long to wait before timeout
* `group_type` - (Required) Defines a type of groups used: static (a group contains a list of users) or dynamic (a user contains a list of groups). Required if ldapGroupsAsRoles is true.
* `host` - (Required) LDAP server connection hostname
* `max_incident_count` - (Required) How many retry attempts
* `name` - (Required) LDAP server name
* `port` - (Required) LDAP server connection port to use
* `protocol` - (Required) LDAP server connection Protocol to use
* `search_base` - (Required) DAP location to be added to the connection URL
* `auth_password` - (Optional) The password to bind with. Required if authScheme other than none.
* `auth_realm` - (Optional) The SASL realm to bind to. Required if authScheme is CRAM_MD5 or DIGEST_MD5
* `group_base_dn` - (Optional) The relative DN where group objects are found (e.g. ou=Group). This value will have the Search base DN value appended to form the full Group search base DN.
* `group_id_attribute` - (Optional) This field specifies the attribute of the Object class that defines the Group ID. Required if groupType is static
* `group_member_attribute` - (Optional) LDAP attribute containing the usernames for the group. Required if groupType is static
* `group_member_format` - (Optional) The format of user ID stored in the group member attribute. Required if groupType is static
* `group_object_class` - (Optional) LDAP class for group objects. Required if groupType is static
* `group_subtree` - (Optional) Are groups located in structures below the group base DN
* `ldap_groups_as_roles` - (Optional) Denotes whether LDAP assigned roles are used as Nexus Repository Manager roles
* `use_trust_store` - (Optional) Whether to use certificates stored in Nexus Repository Manager's truststore
* `user_base_dn` - (Optional) The relative DN where user objects are found (e.g. ou=people). This value will have the Search base DN value appended to form the full User search base DN.
* `user_email_address_attribute` - (Optional) This is used to find an email address given the user ID
* `user_id_attribute` - (Optional) This is used to find a user given its user ID
* `user_ldap_filter` - (Optional) LDAP search filter to limit user search
* `user_member_of_attribute` - (Optional) Set this to the attribute used to store the attribute which holds groups DN in the user object. Required if groupType is dynamic
* `user_object_class` - (Optional) LDAP class for user objects
* `user_password_attribute` - (Optional) If this field is blank the user will be authenticated against a bind with the LDAP server
* `user_real_name_attribute` - (Optional) This is used to find a real name given the user ID
* `user_subtree` - (Optional) Are users located in structures below the user base DN?


