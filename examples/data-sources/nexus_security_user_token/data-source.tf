data "nexus_security_user_token" "nexus" {}

output "nexus_user_token_enabled" {
  description = "User Tokens enabled?"
  value       = data.nexus_security_user_token.nexus.enabled
}
