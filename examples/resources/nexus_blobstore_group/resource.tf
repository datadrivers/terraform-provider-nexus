resource "nexus_blobstore_group" "example" {
  name        = "group-example"
  fill_policy = "roundRobin"
  members = [
    nexus_blobstore_file.one.name,
    nexus_blobstore_file.two.name
  ]
}
