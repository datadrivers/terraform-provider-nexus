# Retrieve Cert via Nexus
data "nexus_security_ssl" "ldap_cert" {
  host = "google.de"
  port = 443
}

data "nexus_security_ssl" "bing" {
  host = "bing.com"
  port = 443
}

# Import Cert into Nexus
resource "nexus_security_ssl_truststore" "ldap_cert" {
  pem = data.nexus_security_ssl.ldap_cert.pem
}

# Import 
resource "nexus_security_ssl_truststore" "bing" {
  pem = data.nexus_security_ssl.bing.pem
}

data "nexus_security_ssl_truststore" "test" {
}

data "nexus_security_ssl_truststore" "test2" {
}

output "truststore" {
  value = data.nexus_security_ssl_truststore.test.certificates
}

output "test" {
  value = resource.nexus_security_ssl_truststore.bing
}
