resource "nexus_security_user_roles" "anonymous" {
  userid    = "anonymous"
  roles     = ["nx-anonymous", "example-role"]
}
