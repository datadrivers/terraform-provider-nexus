resource "nexus_repository_apt_hosted" "bullseye_stable" {
  name   = "bullseye-stable"
  online = true

  distribution = "bullseye"
  signing {
    keypair    = "keypair"
    passphrase = "passphrase"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
