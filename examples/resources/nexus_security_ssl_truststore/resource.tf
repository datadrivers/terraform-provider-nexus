resource "nexus_security_ssl_truststore" "ldap_cert" {
  pem = file("${path.module}/cert.pem")
}
