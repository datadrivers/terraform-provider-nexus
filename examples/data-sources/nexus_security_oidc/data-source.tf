data "nexus_security_oidc" "oidc" {}

output "oidc_authorization_url" {
  description = "OpenID Provider authorization endpoint configured on Nexus"
  value       = data.nexus_security_oidc.oidc.authorization_url
}
