# Terraform provider Nexus

## Introduction

Terraform provider to configure Sonatype Nexus using it's API.

Implemented and tested with Sonatype Nexus 3.20.1.

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
  status    = "online"
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

## Author

[Datadrivers GmbH](https://www.datadrivers.de)