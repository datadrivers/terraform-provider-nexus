resource "nexus_security_realms" "example" {
  active = [
    "NexusAuthenticatingRealm",
    "NexusAuthorizingRealm",
  ]
}
