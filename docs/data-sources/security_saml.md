---
layout: "nexus"
page_title: "Nexus: nexus_security_saml"
sidebar_current: "docs-nexus-datasource-security_saml"
description: |-
  Use this data source to get the saml configuration
---

# nexus_security_saml

Use this data source to get the saml configuration

---
**PRO Feature**
---

## Example Usage

```hcl
data "nexus_security_saml" "saml" {}

output "saml_entity_id" {
  description = "Entity ID URI of saml config"
  value       = data.nexus_security_saml.saml.entity_id
}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `email_attribute` - IdP field mappings for user's email address
* `entity_id` - Entity ID URI
* `first_name_attribute` - IdP field mappings for user's given name
* `groups_attribute` - IdP field mappings for user's groups
* `idp_metadata` - SAML Identity Provider Metadata XML
* `last_name_attribute` - IdP field mappings for user's family name
* `username_attribute` - IdP field mappings for username
* `validate_assertion_signature` - By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the assertions.
* `validate_response_signature` - By default, if a signing key is found in the IdP metadata, then NXRM will attempt to validate signatures on the response.


