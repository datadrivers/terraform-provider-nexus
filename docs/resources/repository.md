---
layout: "nexus"
page_title: "Nexus: nexus_repository"
subcategory: "Other"
sidebar_current: "docs-nexus-resource-repository"
description: |-
  Use this resource to create a Nexus Repository.
---

# nexus_repository

Use this resource to create a Nexus Repository.

## Example Usage

Example Usage - Apt hosted repository

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

Example Usage - Docker group repository

```hcl
resource "nexus_repository" "docker_group" {
  name   = "docker-group"
  format = "docker"
  type   = "group"
  online = true

  group {
    member_names = [
      "docker_releases",
      "docker_hub"
    ]
  }

  docker {
    force_basic_auth = false
    http_port        = 5000
    https_port       = 0
    v1enabled        = false
  }

  storage {
    blob_store_name                = "docker_group_blobstore"
    strict_content_type_validation = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `format` - (Required, ForceNew) Repository format. Possible values: `apt`, `bower`, `conan`, `docker`, `gitlfs`, `go`, `helm`, `maven2`, `npm`, `nuget`, `p2`, `pypi`, `raw`, `rubygems`, `yum`
* `name` - (Required) A unique identifier for this repository
* `type` - (Required, ForceNew) Repository type. Possible values: `group`, `hosted`, `proxy`
* `apt_signing` - (Optional) Apt signing configuration for the repository
* `apt` - (Optional) Apt specific configuration of the repository
* `bower` - (Optional) Bower specific configuration of the repository
* `cleanup` - (Optional) Cleanup policies
* `docker_proxy` - (Optional) Configuration for docker proxy repository
* `docker` - (Optional) Docker specific configuration of the repository
* `group` - (Optional) Configuration for repository group
* `http_client` - (Optional) HTTP Client configuration for proxy repositories
* `maven` - (Optional) Maven specific configuration of the repository
* `negative_cache` - (Optional) Configuration of the negative cache handling
* `nuget_proxy` - (Optional) Configuration for the nuget proxy repository
* `online` - (Optional) Whether this repository accepts incoming requests
* `proxy` - (Optional) Configuration for the proxy repository
* `storage` - (Optional) The storage configuration of the repository
* `yum` - (Optional) Yum specific configuration of the repository

The `apt` object supports the following:

* `distribution` - (Required) The linux distribution
* `flat` - (Optional) Whether this repository is flat

The `apt_signing` object supports the following:

* `keypair` - (Required) PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor )
* `passphrase` - (Required) Passphrase for the keypair

The `bower` object supports the following:

* `rewrite_package_urls` - (Optional) Force Bower to retrieve packages through the proxy repository

The `cleanup` object supports the following:

* `policy_names` - (Optional) List of policy names

The `docker` object supports the following:

* `force_basic_auth` - (Optional) Whether to force authentication (Docker Bearer Token Realm required if false)
* `http_port` - (Optional) Create an HTTP connector at specified port
* `https_port` - (Optional) Create an HTTPS connector at specified port
* `v1enabled` - (Optional) Whether to allow clients to use the V1 API to interact with this repository

The `docker_proxy` object supports the following:

* `index_type` - (Required) Type of Docker Index
* `index_url` - (Optional) URL of Docker Index to use

The `group` object supports the following:

* `member_names` - (Required) Member repositories names

The `http_client` object supports the following:

* `authentication` - (Optional) Authentication configuration of the HTTP client
* `auto_block` - (Optional) Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive
* `blocked` - (Optional) Whether to block outbound connections on the repository
* `connection` - (Optional) Connection configuration of the HTTP client

The `authentication` object supports the following:

* `ntlm_domain` - (Optional) The ntlm domain to connect
* `ntlm_host` - (Optional) The ntlm host to connect
* `password` - (Optional) The password used by the proxy repository
* `type` - (Optional) Authentication type. Possible values: `ntlm`, `username` or `bearerToken`. Only npm supports bearerToken authentication
* `username` - (Optional) The username used by the proxy repository

The `connection` object supports the following:

* `enable_cookies` - (Optional) Whether to allow cookies to be stored and used
* `retries` - (Optional) Total retries if the initial connection attempt suffers a timeout
* `timeout` - (Optional) Seconds to wait for activity before stopping and retrying the connection
* `use_trust_store` - (Optional) Use certificates stored in the Nexus Repository Manager truststore to connect to external systems
* `user_agent_suffix` - (Optional) Custom fragment to append to User-Agent header in HTTP requests

The `maven` object supports the following:

* `layout_policy` - (Optional) Validate that all paths are maven artifact or metadata paths. Possible values: `PERMISSIVE` or `STRICT`
* `version_policy` - (Optional) What type of artifacts does this repository store? Possible values: `RELEASE`, `SNAPSHOT` or `MIXED`

The `negative_cache` object supports the following:

* `enabled` - (Optional) Whether to cache responses for content not present in the proxied repository
* `ttl` - (Optional) How long to cache the fact that a file was not found in the repository (in minutes)

The `nuget_proxy` object supports the following:

* `query_cache_item_max_age` - (Required) What type of artifacts does this repository store
* `nuget_version` - (Optional) Nuget protocol version. Possible values: `V2` or `V3` (Default)

The `proxy` object supports the following:

* `content_max_age` - (Optional) How long (in minutes) to cache artifacts before rechecking the remote repository
* `metadata_max_age` - (Optional) How long (in minutes) to cache metadata before rechecking the remote repository.
* `remote_url` - (Optional) Location of the remote repository being proxied

The `storage` object supports the following:

* `blob_store_name` - (Optional) Blob store used to store repository contents
* `strict_content_type_validation` - (Optional) Whether to validate uploaded content's MIME type appropriate for the repository format
* `write_policy` - (Optional) Controls if deployments of and updates to assets are allowed. Possible values: `ALLOW`, `ALLOW_ONCE`, `DENY`

The `yum` object supports the following:

* `deploy_policy` - (Required) Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`
* `repodata_depth` - (Optional) Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5


