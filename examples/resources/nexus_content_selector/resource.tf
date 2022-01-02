resource "nexus_content_selector" "docker_public" {
  name        = "docker-public"
  description = "A content selector matching public docker images."
  expression  = "path =^ \"/v2/public/\""
}
