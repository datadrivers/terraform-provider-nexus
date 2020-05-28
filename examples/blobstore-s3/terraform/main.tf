// iam user
resource "aws_iam_user" "login_for_external_user" {
  count = var.iam_user_name_prefix != null ? 1 : 0
  name  = var.iam_user_name_prefix
  tags  = var.tags
}

resource "aws_iam_user_policy" "assume_nexus_role" {
  count = var.iam_user_name_prefix != null ? 1 : 0
  name  = "assume_role_nexus"
  user  = aws_iam_user.login_for_external_user.0.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Resource": "${aws_iam_role.access_s3_bucket_nexus.arn}"
    }
  ]
}
EOF
}

// iam role
resource "aws_iam_role" "access_s3_bucket_nexus" {
  name = var.iam_role_name

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "AWS": ${jsonencode(aws_iam_user.login_for_external_user.*.arn)}
      },
      "Effect": "Allow"
    }
  ]
}
EOF

  tags = var.tags
}

# ref: https://help.sonatype.com/repomanager3/high-availability/configuring-blob-stores#ConfiguringBlobStores-AWSSimpleStorageService(S3)
data "aws_iam_policy_document" "bucket_policy" {
  statement {
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject",
      "s3:ListBucket",
      "s3:GetLifecycleConfiguration",
      "s3:PutLifecycleConfiguration",
      "s3:PutObjectTagging",
      "s3:GetObjectTagging",
      "s3:DeleteObjectTagging",
      "s3:DeleteBucket",
      "s3:CreateBucket",
      "s3:GetBucketAcl"
    ]

    resources = [
      module.s3_bucket.this_s3_bucket_arn,
      format("%s/*", module.s3_bucket.this_s3_bucket_arn),
    ]
  }
}

resource "aws_iam_policy" "policy_s3_bucket" {
  name_prefix = var.iam_role_name
  description = "Nexus Access to S3 bucket"

  policy = data.aws_iam_policy_document.bucket_policy.json
}

resource "aws_iam_role_policy_attachment" "attach-s3-policy" {
  role       = aws_iam_role.access_s3_bucket_nexus.name
  policy_arn = aws_iam_policy.policy_s3_bucket.arn
}

// s3 bucket
module "s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket              = var.bucket_name
  bucket_prefix       = var.bucket_prefix
  acl                 = var.acl
  force_destroy       = var.force_destroy
  acceleration_status = var.acceleration_status
  region              = var.region
  request_payer       = var.request_payer

  cors_rule                            = var.cors_rule
  lifecycle_rule                       = var.lifecycle_rule
  logging                              = var.logging
  object_lock_configuration            = var.object_lock_configuration
  replication_configuration            = var.replication_configuration
  server_side_encryption_configuration = var.server_side_encryption_configuration
  versioning                           = var.versioning
  website                              = var.website
  block_public_acls                    = var.block_public_acls
  block_public_policy                  = var.block_public_policy
  ignore_public_acls                   = var.ignore_public_acls
  restrict_public_buckets              = var.restrict_public_buckets

  tags = var.tags
}
