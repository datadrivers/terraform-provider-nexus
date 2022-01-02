resource "nexus_security_content_selector" "example" {
	name        = "example"
	description = "example content selector"
	expression  = "format == \"raw\""
}