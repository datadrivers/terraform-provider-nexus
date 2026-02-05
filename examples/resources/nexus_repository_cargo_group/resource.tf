resource "nexus_repository_cargo_hosted" "releases" {
  name   = "cargo-releases"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }

}


resource "nexus_repository_cargo_group" "group" {
  name   = "cargo-group"
  online = true

  group {
    member_names = [
      nexus_repository_cargo_hosted.releases.name,
    ]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
