# Retrieve certificates from Nexus truststore
data "nexus_security_ssl_truststore" "nexus_truststore" {
}

# Output Nexus truststore certificates
output "truststore" {
  value = data.nexus_security_ssl_truststore.nexus_truststore
}
