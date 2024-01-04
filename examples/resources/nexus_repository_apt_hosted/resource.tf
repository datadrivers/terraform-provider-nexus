resource "nexus_repository_apt_hosted" "bullseye_stable" {
  name   = "bullseye-stable"
  online = true

  distribution = "bullseye"

  signing {
    # The passphrase set here will never be returned from the nexus API.
    # When reading the resource, the passphrase will be read from the previous state,
    # so external changes won't be detected.`,

    # If the passphrase is unset or empty, the nexus API will also not return the keypair.
    # When reading the resource, the keypair will be read from the previous state in this case,
    # so external changes won't be detected.`,

    keypair    = "keypair"
    passphrase = "passphrase"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }
}
