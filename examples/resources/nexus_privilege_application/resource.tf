resource "nexus_privilege_application" "example" {
  name = "example_privilege"
  description = "description"
  actions = [ "ADD", "READ", "EDIT", "DELETE" ]
  domain = "domain"
}
