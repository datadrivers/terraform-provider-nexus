---
page_title: "Resource nexus_security_ssl_truststore"
subcategory: "Security"
description: |-
  Use this resource to add an SSL certificate to the nexus Truststore
---
# Resource nexus_security_ssl_truststore
Use this resource to add an SSL certificate to the nexus Truststore
## Example Usage
```terraform
resource "nexus_security_ssl_truststore" "ldap_cert" {
  pem = file("${path.module}/cert.pem")
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `pem` (String) The cert in PEM format

### Read-Only

- `fingerprint` (String) The fingerprint of the cert
- `id` (String) Used to identify resource at nexus
