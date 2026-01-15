# Configure Nexus IQ Server integration
resource "nexus_iqserver" "main" {
  enabled             = true
  url                 = "https://iq-server.example.com"
  authentication_type = "USER"
  username            = "admin"
  password            = "admin123"

  show_link               = false
  use_trust_store_for_url = true
  fail_open_mode_enabled  = false
  timeout_seconds         = 60
}

# Example: IQ Server with PKI authentication
resource "nexus_iqserver" "pki" {
  enabled             = true
  url                 = "https://iq-server.example.com"
  authentication_type = "PKI"

  show_link               = true
  use_trust_store_for_url = true
  fail_open_mode_enabled  = false
  timeout_seconds         = 30
}
