resource "nexus_blobstore_azure" "example" {
  name = "example"

  bucket_configuration {
    account_name = "example-account-name"
    authentication {
      authentication_method = "ACCOUNTKEY"
      account_key           = "example-account-key"
    }
    container_name = "example-container-name"
  }

  soft_quota {
    limit = 1024
    type  = "spaceRemainingQuota"
  }
}
