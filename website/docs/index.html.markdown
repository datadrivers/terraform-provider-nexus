---
layout: "nexus"
page_title: "Provider: Nexus"
sidebar_current: "docs-nexus-index"
description: |-
  The Nexus provider allows Terraform to read from, write to, and configure Sonatype Nexus Repository Manager
---

# Nexus Provider
  
The Nexus provider allows Terraform to read from, write to, and configure [Sonatype Nexus Repository Manager](https://de.sonatype.com/product-nexus-repository).

-> **Note** This provider hat been implemented and tested with Sonatype Nexus `3.25.1-04`.

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

### Data Sources

#### Data nexus_blobstore

```hcl
data "nexus_blobstore" "default" {
  name = "default
}
```

#### Data nexus_privileges

Get all privileges matching all optional filters. All parameters are optional.
The returned list contains all privileges that match all specified parameters!

```hcl
data "nexus_privileges" "exmaple" {
  domain     = "application"
  format     = "maven2"
  repository = "maven-public"
  type       = "repository-admin"
}
```

#### Data nexus_repository

```hcl
data "nexus_repository" "maven-central" {
  name = "maven-central"
}
```

#### Data nexus_security_ldap

Return LDAP server

```hcl
data "nexus_security_ldap" "example" {}
```

#### Data nexus_security_realms

Return active and available security realms

```hcl
data "nexus_security_realms" "example" {}
```

#### Data nexus_user

```hcl
data "nexus_user" "admin" {
  userid = "admin"
}
```

### Resources

#### Resource nexus_blobstore

Blobstore can be imported using

```shell
terraform import nexus_blobstore.default default
```

##### Use File

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

##### Use S3

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

#### Resource nexus_content_selector

Content selector can be imported using

```shell
terraform import nexus_content_selector.docker-public docker-public
```

```hcl
resource "nexus_content_selector" "docker-public" {
  name = "docker-public"
  description = "A content selector matching public docker images."
  expression = "path =^ \"/v2/public/\""
}
```

#### Resource nexus_privilege

Privilege can be imported using

```shell
terraform import nexus_privilege.docker-public-read docker-public-read
```

##### Content Selector

```hcl
resource "nexus_privilege" "docker-public-read" {
  name = "docker-public-read"
  description = "Read permission on the docker public path."
  type = "repository-content-selector"
  content_selector = "docker-public"
  actions = [
    "read"
  ]
}
```

#### Resource nexus_repository

Repository can be imported using

```shell
terraform import nexus_repository.maven_central maven-central
```

##### APT hosted

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

##### Bower hosted

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

##### Docker group

```hcl
resource "nexus_repository" "docker_group" {
  name   = "docker-group"
  format = "docker"
  type   = "group"
  online = true

  group {
    member_names = ["docker-hub"]
  }

  docker {
    force_basic_auth = true
    http_port        = 5000
    https_port       = 5001
    v1enabled        = false
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
```

##### Docker hosted

```hcl
resource "nexus_repository" "docker_hosted" {
  name   = "docker-hosted"
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
resource "nexus_repository" "docker_hub" {
  name   = "docker-hub"
  type   = "proxy"
  format = "docker"

  docker {
    force_basic_auth = true
    v1enabled        = true
  }

  docker_proxy {
    index_type = "HUB"
  }

  http_client {
    # Optional
    authentication {
      type        = "username"
      username    = "example"
      ntlm_domain = "example"
      ntlm_host   = "host.example.com"
    }
  }

  negative_cache {
    enabled = true
    ttl     = 1440
  }

  proxy {
    remote_url  = "https://registry-1.docker.io"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
  }
}
```

##### NPM hosted

```hcl
resource "nexus_repository" "npm_hosted" {
  name   = "npm-hosted-repo"
  format = "npm"
  type   = "hosted"

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

##### Nuget proxy

```hcl
resource "nexus_repository" "nuget_proxy" {
  name   = "nuget-proxy-repo"
  format = "nuget"
  type   = "proxy"
  online = true

  http_client {
    auto_block = true
    blocked    = false

    # Optional
    authentication {
      type        = "username"
      username    = "example"
      ntlm_domain = "example"
      ntlm_host   = "host.example.com"
    }
}

    negative_cache {
	enabled = true
	ttl     = 1440
    }

    nuget_proxy {
	query_cache_item_max_age = 1440
    }

    proxy {
	remote_url  = "https://www.nuget.org/api/v2/"
    }

    storage {
	write_policy = "ALLOW"
    }
}
```

##### PyPi hosted

```hcl
resource "nexus_repository" "pypi_hosted" {
  name   = "pypi-hosted-repo"
  format = "pypi"
  type   = "hosted"

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

##### PyPi proxy

```hcl
resource "nexus_repository" "pypi_proxy" {
  name   = "pypi-proxy-repo"
  format = "pypi"
  type   = "proxy"
  online = true

  http_client {
    auto_block = true
    blocked    = false

    # Optional
    authentication {
      type        = "username"
      username    = "example"
      ntlm_domain = "example"
      ntlm_host   = "host.example.com"
    }
}

    negative_cache {
	enabled = true
	ttl     = 1440
    }

    proxy {
	remote_url  = "https://pypi.org"
    }

    storage {
	write_policy = "ALLOW"
    }
}
```

##### RAW hosted

```hcl
resource "nexus_repository" "raw_hosted" {
  name   = "raw-hosted-repo"
  format = "raw"
  type   = "hosted"

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

##### RAW proxy

```hcl
resource "nexus_repository" "raw_proxy" {
    format = "raw"
    name   = "raw-proxy-repo"
    online = true
    type   = "proxy"

    proxy {
	remote_url  = "https://nodejs.org/dist/"
    }

    http_client {
    ...
    }

    negative_cache {
	enabled = true
	ttl     = 1440
    }

    storage {
    ...
    }
}
```

#### Resource nexus_role

Role can be imported using

```shell
terraform import nexus_role.nx_admin nx-admin
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

#### Resource nexus_security_ldap

Configure LDAP server

```shell
terraform import nexus_security_ldap.example example
```

```hcl
resource "nexus_security_ldap" "acceptance" {
  auth_password                  = "t0ps3cr3t"
  auth_realm                     = "EXAMPLE"
  auth_schema                    = ""
  auth_username                  = "admin"
  connection_retry_delay_seconds = 1
  connection_timeout_seconds     = 1
  group_base_dn                  = "ou=Group"
  group_id_attribute             = "cn"
  group_member_attribute         = "memberUid"
  group_member_format            = "uid=${username},ou=people,dc=example,dc=com"
  group_object_class             = "example"
  group_subtree                  = true
  host                           = "ldap.example.com"
  ldap_groups_as_roles           = true
  max_incident_count             = 1
  name                           = "example-ldap"
  port                           = 389
  protocol                       = "LDAP"
  search_base                    = "dc=example,dc=com"
  user_base_dn                   = "ou=people"
  user_email_address_attribute   = "mail"
  user_id_attribute              = "uid"
  user_ldap_filter               = "(|(mail=*@example.com)(uid=dom*))"
  user_member_of_attribute       = "memberOf"
  user_object_class              = "posixGroup"
  user_password_attribute        = "exmaple"
  user_real_name_attribute       = "cn"
  user_subtree                   = true
}
```

#### Resource nexus_security_ldap_order

Set order of LDAP server

```hcl
resource "nexus_security_ldap" "server1" {
  ...
  name = "server1"
}

resource "nexus_security_ldap" "server2" {
  ...
  name = "server2"
}

resource "nexus_security_ldap_order" {
  order = [
    nexus_security_ldap.server1.name,
    nexus_security_ldap.server2.name,
  ]
}
```

#### Resource nexus_security_realms

Activate security realms

```shell
terraform import nexus_security_realms.example
```

```hcl
resource "nexus_security_realms" "example" {
  active = ["NexusAuthenticatingRealm", "NexusAuthorizingRealm"]
}
```

#### Resource nexus_script

Script can be imported using

```shell
terraform import nexus_script.my_script my-script
```

```hcl
resource "nexus_script" "hello_world" {
  name    = "hello-world"
  content = "log.info('Hello, World!')"
}
```

#### Resource nexus_user

User can be imported using

```shell
terraform import nexus_user.admin admin
```

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

## Author

[Datadrivers GmbH](https://www.datadrivers.de)
