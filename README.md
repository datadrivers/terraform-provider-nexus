# Terraform provider Nexus

## Introduction

Terraform provider to configure Sonatype Nexus using it's API.

Implemented and tested with Sonatype Nexus 3.21.2.

## Usage

### Provider config

```hcl
provider "nexus" {
  url      = "http://127.0.0.1:8080"
  username = "admin"
  password = "admin123"
}
```

### Data

#### nexus_blobstore

```hcl
data "nexus_blobstore" "default" {
  name = "default
}
```

#### nexus_repository

```hcl
data "nexus_repository" "maven-central" {
  name = "maven-central"
}
```

#### nexus_user

```hcl
data "nexus_user" "admin" {
  userid = "admin"
}
```

### Resources

#### nexus_blobstore

Blobstore can be imported using

```shell
$ terraform import nexus_blobstore.default default
```

##### File

```hcl
resource "nexus_blobstore" "default" {
  name = "blobstore-01"
  type = "File"
  path = "/nexus-data/blobstore-01"

  soft_quota {
    limit = 1024
    type  = "spaceRemainingQuota"
  }
}
```

##### S3

```hcl
resource "nexus_blobstore" "aws" {
  name = "blobstore-01"
  type = "S3"

  bucket_configuration {
    bucket {
      name   = "aws-bucket-name"
      region = "us-central-1"
    }

    bucket_security {
      access_key_id = "<your-aws-access-key-id>"
      secret_access_key = "<your-aws-secret-access-key>"
    }
  }

  soft_quota {
    limit = 1024
    type  = "spaceRemainingQuota"
  }
}
```

#### nexus_repository

Repository can be imported using
```shell
$ terraform import nexus_repository.maven_central maven-central
```

```hcl
resource "nexus_repository" "apt_hosted" {
  name   = "apt-repo"
  format = "apt"
  type   = "hosted"

  apt {
    distribution = "bionic"
  }

  apt_signing {
    keypair    = "<keypair>"
    passphrase = "<passphrase>"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

```hcl
resource "nexus_repository" "bower_hosted" {
  name   = "bower-hosted-repo"
  format = "bower"
  type   = "hosted"

  bower {
    rewrite_package_urls = false
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

```hcl
resource "nexus_repository" "docker_hosted" {
  name   = "docker-hosted-repo"
  format = "docker"
  type   = "hosted"
  online = true

  docker {
    http_port        = 8082
    https_port       = 8083
    force_basic_auth = true
    v1enabled        = true
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

```hcl
resource "nexus_repository" "docker_proxy" {
  name   = "docker-proxy-repo"
  type   = "proxy"
  format = "docker"

  docker {
    http_port        = 8082
    https_port       = 8083
    force_basic_auth = true
    v1enabled        = true
  }

  docker_proxy {
    index_url  = "https://index.docker.io"
    index_type = "HUB"
  }

  http_client {

  }

  negative_cache {

  }

  proxy {

  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

#### nexus_role

Role can be imported using
```shell
$ terraform import nexus_role.nx_admin nx-admin
```

```hcl
resource "nexus_role" "nx-admin" {
  roleid      = "nx-admin"
  name        = "nx-admin"
  description = "Administrator role"
  privileges  = ["nx-all"]
  roles       = []
}
```

#### nexus_user

User can be imported using
```shell
$ terraform import nexus_user.admin admin
````

```hcl
resource "nexus_user" "admin" {
  userid    = "admin"
  firstname = "Administrator"
  lastname  = "User"
  email     = "nexus@example.com"
  password  = "admin123"
  roles     = ["nx-admin"]
  status    = "active"
}
```

#### nexus_script

Script can be imported using
```shell
$ terraform import nexus_script.my_script my-script
```

```hcl
resource "nexus_script" "hello_world" {
  name    = "hello-world"
  content = "log.info('Hello, World!')"
}
```

## Build

There is a [makefile](./GNUmakefile) to build the provider.

```sh
make
```

## Testing

For testing start a local Docker container using script [./scripts/start-nexus.sh](./scripts/start-nexus.sh).

```shell
$ ./scripts/start-nexus.sh
```

This will start a Docker container and expose port 8081.

Now start the tests

```shell
$ NEXUS_URL="http://127.0.0.1:8081" NEXUS_USERNAME="admin" NEXUS_PASSWORD="admin123" make testacc
```

__NOTE__: To test Blobstore type S3 following environment variables must be set, otherwise tests will fail

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_DEFAULT_REGION` the AWS region of the S3 bucket to use, defaults to `eu-central-1`
- `AWS_BUCKET_NAME` the name of S3 bucket to use, defaults to `terraform-provider-nexus-s3-test`

## Author

[Datadrivers GmbH](https://www.datadrivers.de)
