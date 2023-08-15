resource "nexus_privilege_application" "priv" {
  name = "new-app-privilege"
  description = "description"
  actions = ["UPDATE"]
  domain = "domain"
}
