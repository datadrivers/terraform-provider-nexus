data "nexus_security_anonymous" "nexus" {
}

output "nexus_anonymous_enabled" {
  description = "Anonymous enabled?"
  value       = data.nexus_security_anonymous.nexus.enabled
}
