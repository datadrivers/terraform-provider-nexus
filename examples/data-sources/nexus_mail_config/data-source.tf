
data "nexus_mail_config" "config" {}
output "host" {
  value = data.nexus_mail_config.config.host
}
