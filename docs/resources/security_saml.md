---
layout: "nexus"
page_title: "Nexus: nexus_security_saml"
subcategory: "Security"
sidebar_current: "docs-nexus-resource-security_saml"
description: |-
  Use this resource to create a Nexus Security SAML
---

# nexus_security_saml

Use this resource to create a Nexus Security SAML

---
**PRO Feature**
---

## Example Usage

```hcl
resource "nexus_security_saml" "example" {
  idp_metadata                  = "<EntityDescriptor ...>...</EntityDescriptor>"
  entity_id                     = "http://nexus.example/service/rest/v1/security/saml/metadata"
  validate_response_signature   = true
  validate_assertion_signature  = true
  username_attribute            = "username"
  first_name_attribute          = "firstName"
  last_name_attribute           = "lastName
  email_attribute               = "email
  groups_attribute              = "groups"
}
```

## Argument Reference

The following arguments are supported:

* `idp_metadata` - (Required) SAML Identity Provider Metadata XML
* `username_attribute` - (Required) IdP field mappings for username
* `email_attribute` - (Optional) IdP field mappings for user's email address
* `entity_id` - (Optional) Entity ID URI
* `first_name_attribute` - (Optional) IdP field mappings for user's given name
* `groups_attribute` - (Optional) IdP field mappings for user's groups
* `last_name_attribute` - (Optional) IdP field mappings for user's family name
* `validate_assertion_signature` - (Optional) By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the assertions.
* `validate_response_signature` - (Optional) By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the response.


