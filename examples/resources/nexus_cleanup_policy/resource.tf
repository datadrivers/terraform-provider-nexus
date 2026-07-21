package nexus_cleanup_policy
resource "nexus_cleanup_policy" "example" {
name                      = "maven-artifacts-clean-up"
notes                     = "This is an example cleanup policy"
criteria_last_downloaded   = 0
criteria_release_type      = "RELEASES"
criteria_asset_regex       = ".*"
retain                    = null
format                    = "maven2"
}