data "nexus_anonymous" "nexus" {
}

output "nexus_anonymous_enabled" {
  description = "Anonymous enabled?"
  value       = data.nexus_anonymous.nexus.enabled
}
