bucket_name     = "terrafor-provider-nexus--example"
iam_name_prefix = "terrafor-provider-nexus"

acl = "private"

force_destroy = true

versioning = {
  enabled = false
}

lifecycle_rule = [
  {
    id      = "save-costs"
    enabled = true

    abort_incomplete_multipart_upload_days = 7

    transition = [
      {
        days          = 30
        storage_class = "STANDARD_IA" # (requires >=30 days)
      },
    ]
  },
]

server_side_encryption_configuration = {
  rule = {
    apply_server_side_encryption_by_default = {
      sse_algorithm = "AES256"
    }
  }
}

tags = {
  provisioning = "terraform"
  use-case     = "terrafor-provider-nexus"
  owner        = "datadrivers"
}

// S3 bucket-level Public Access Block configuration
block_public_acls = true

block_public_policy = true

ignore_public_acls = true

restrict_public_buckets = true
