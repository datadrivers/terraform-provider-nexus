resource "nexus_repository_maven_hosted" "releases" {
  name   = "maven-releases"
  online = true

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = false
    write_policy                   = "ALLOW"
  }

  maven {
    version_policy      = "RELEASE"
    layout_policy       = "STRICT"
    content_disposition = "INLINE"
  }
}


resource "nexus_repository_maven_group" "group" {
  name   = "maven-group"
  online = true

  group {
    member_names = [
      nexus_repository_maven_hosted.releases.name,
    ]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
