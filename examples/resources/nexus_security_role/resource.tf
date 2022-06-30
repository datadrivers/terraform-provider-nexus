# Example Usage - Create a group with roles
resource "nexus_security_role" "nx_admin" {
  roleid      = "nx-admin"
  name        = "nx-admin"
  description = "Administrator role"
  privileges  = ["nx-all"]
  roles       = []
}

# Example Usage - Create a group with privileges
resource "nexus_security_role" "docker_deploy" {
  description = "Docker deployment role"
  name        = "docker-deploy"
  privileges = [
    "nx-repository-view-docker-*-*",
  ]
  roleid = "docker-deploy"
}
