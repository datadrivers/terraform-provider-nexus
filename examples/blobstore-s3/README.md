# Nexus with AWS s3 blobstore

To use the AWS credential from the terraform code example below you need to do the following:

- login with iam user credentials
- assume role "terraform-provider-nexus-s3"
- access s3 bucket "terraform-provider-nexus-example.datadrivers.de"

## Terraform

Here you can find terraform example code to create iam user, iam role and s3-bucket to be used with nexus.

```bash
cd terraform
terraform init
terraform plan -var-file=./example.tfvars -out example.tfplan
terraform apply example.tfvars
