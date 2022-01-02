data "nexus_security_saml" "saml" {}

output "saml_entity_id" {
  description = "Entity ID URI of saml config"
  value       = data.nexus_security_saml.saml.entity_id
}
