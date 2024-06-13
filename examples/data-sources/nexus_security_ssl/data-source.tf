# Retrieve Cert via Nexus
data "nexus_security_ssl" "ldap_cert" {
  host = "google.de"
  port = 443
}

# Import Cert into Nexus
resource "nexus_security_ssl_truststore" "ldap_cert" {
  pem = data.nexus_security_ssl.ldap_cert.pem
}
