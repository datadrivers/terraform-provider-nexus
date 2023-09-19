resource "nexus_privilege_script" "example_privilege_script" {
	name = "new-script-privilege"
	description = "description"
	actions = ["ADD", "READ", "DELETE", "RUN", "BROWSE", "EDIT"]
	script_name = resource.nexus_script.some_script.name
}