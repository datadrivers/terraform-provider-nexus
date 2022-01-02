resource "nexus_security_ldap" "server1" {
  ...
  name = "server1"
}

resource "nexus_security_ldap" "server2" {
  ...
  name = "server2"
}

resource "nexus_security_ldap_order" "system" {
  order = [
    nexus_security_ldap.server1.name,
    nexus_security_ldap.server2.name,
  ]
}
