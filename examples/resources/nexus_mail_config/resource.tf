resource "nexus_mail_config" "config" {
  port         = 25
  host         = "examplehost.org"
  from_address = "from@examplehost.org"

  enabled                   = true
  username                  = "uname"
  password                  = "topsecret"
  subject_prefix            = "prefix: "
  start_tls_enabled         = true
  start_tls_required        = true
  ssl_on_connect_enabled    = true
  nexus_trust_store_enabled = true
}
