data "nexus_privileges" "example" {
  domain     = "application"
  format     = "maven2"
  repository = "maven-public"
  type       = "repository-admin"
}
