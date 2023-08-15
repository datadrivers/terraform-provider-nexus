data "nexus_privilege_application" "priv" {
  name = "privilege-name"
}

output "privilege_name" {
  value = data.nexus_privilege_application.priv.name
}
