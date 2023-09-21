resource "nexus_script" "example" {
  name    = "example_script"
  content = "log.info('Hello, World!')"
}

resource "nexus_privilege_script" "example" {
  name        = "example_privilege"
  description = "description"
  actions     = ["ADD", "READ", "DELETE", "RUN", "BROWSE", "EDIT"]
  script_name = resource.nexus_script.example.name
}
