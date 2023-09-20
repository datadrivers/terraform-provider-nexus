resource "nexus_privilege_wildcard" "example" {
  name        = "example_privilege"
  description = "description"
  pattern     = "nexus:*"
}
