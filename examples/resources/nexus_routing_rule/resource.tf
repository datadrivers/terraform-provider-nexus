resource "nexus_routing_rule" "stop_leaks" {
  name        = "stop-leaks"
  description = "Prevent requests of internal names"
  mode        = "BLOCK"
  matchers = [
    "^/com/example/.*",
    "^/org/example/.*",
  ]
}
