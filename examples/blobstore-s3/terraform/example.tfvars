// general
tags = {
  provisioning = "terraform"
  use-case     = "terraform-provider-nexus"
  github       = "https://github.com/dre2004/terraform-provider-nexus"
  owner        = "datadrivers"
}

// iam settings
iam_user_name_prefix = "terraform-provider-nexus--github-action"
iam_role_name        = "terraform-provider-nexus-s3"

// s3 settings
bucket_name   = "terraform-provider-nexus-example.datadrivers.de"
acl           = "private"
force_destroy = true
versioning = {
  enabled = false
}
lifecycle_rule = [] # disable lifecycle rules since it will conflict with nexus lifecycle policys
server_side_encryption_configuration = {
  rule = {
    apply_server_side_encryption_by_default = {
      sse_algorithm = "AES256"
    }
  }
}
// S3 bucket-level Public Access Block configuration
block_public_acls       = true
block_public_policy     = true
ignore_public_acls      = true
restrict_public_buckets = true
