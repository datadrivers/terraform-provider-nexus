# Manage a Nexus capability
# Note: Capabilities API requires Nexus Repository Manager 3.81.0 or later

# Example: Repository Firewall Audit Capability
resource "nexus_capability" "firewall_audit" {
  type    = "firewall.audit"
  enabled = true
  notes   = "Firewall audit capability for maven-central repository"

  properties = {
    repository = "maven-central"
    quarantine = "false"
  }
}

# Example: Outreach Management Capability
resource "nexus_capability" "outreach" {
  type    = "OutreachManagementCapability"
  enabled = true
  notes   = "Enable outreach management"

  properties = {}
}

# Example: HTTP Client Capability
resource "nexus_capability" "http_client" {
  type    = "httpclient"
  enabled = true
  notes   = "Global HTTP client configuration"

  properties = {
    "userAgentCustomisation" = "Nexus Repository Manager"
  }
}
