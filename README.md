Terraform provider Nexus
==========

# Introduction

Terraform provider to configure Sonatype Nexus using it's API.

Implemented and tested with Sonatype Nexus 3.20.1.

# Usage

## Provider config

```hcl
provider "nexus" {
  url      = "http://127.0.0.1:8080"
  username = "admin"
  password = "admin123"
}
```

### Data

#### nexus_role

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

#### nexus_role

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

# Build

There is a [makefile](./GNUmakefile) to build the provider.

```sh
$ make
```

# Author

[Datadrivers GmbH](https://www.datadrivers.de)