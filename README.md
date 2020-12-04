# Terraform provider Nexus

- [Introduction](#introduction)
- [Usage](#usage)
  - [Provider config](#provider-config)
- [Build](#build)
- [Testing](#testing)
- [Create documentation](#create-documentation)
- [Author](#author)

## Introduction

Terraform provider to configure Sonatype Nexus using it's API.

Implemented and tested with Sonatype Nexus `3.29.0-02`.

## Usage

### Provider config

```hcl
provider "nexus" {
  insecure = true
  password = "admin123"
  url      = "https://127.0.0.1:8080"
  username = "admin"
}
```

## Build

There is a [makefile](./GNUmakefile) to build the provider.

```sh
make
```

To build and install provider on macOS into `~/.terraform.d/plugins/darwin_amd64`, you can run

```sh
make darwin-build-install
```

In this case provider will be available to use with your terraform codebase (in terraform init stage).

## Testing

For testing start a local Docker containers using make

```shell
make start-services
```

This will start a Docker and MinIO containers and expose ports 8081 and 9000.

Now start the tests

```shell
NEXUS_URL="http://127.0.0.1:8081" NEXUS_USERNAME="admin" NEXUS_PASSWORD="admin123" AWS_ACCESS_KEY_ID="minioadmin" AWS_SECRET_ACCESS_KEY="minioadmin" AWS_ENDPOINT="http://minio:9000" make testacc
```

or without S3 tests:

```shell
SKIP_S3_TESTS=1 NEXUS_URL="http://127.0.0.1:8081" NEXUS_USERNAME="admin" NEXUS_PASSWORD="admin123" make testacc
```

**NOTE**: To test Blobstore type S3 following environment variables must be set, otherwise tests will fail:

- `AWS_ACCESS_KEY_ID`,
- `AWS_SECRET_ACCESS_KEY`,
- `AWS_DEFAULT_REGION` the AWS region of the S3 bucket to use, defaults to `eu-central-1`,
- `AWS_BUCKET_NAME` the name of S3 bucket to use, defaults to `terraform-provider-nexus-s3-test`.

Optionally you can set `AWS_ENDPOINT` to point an alternative S3 endpoint.

To debug tests

Set env variable `TF_LOG=DEBUG` to see additional output.

Use `printState()` function to discover terraform state (and resource props) during test.

Debug configurations are also available for VS Code.

## Create documentation

To generate the terraform documentation from go files, you can run

```shell
make docs
```

## Author

[Datadrivers GmbH](https://www.datadrivers.de)
