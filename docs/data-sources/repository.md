---
layout: "nexus"
page_title: "Nexus: nexus_repository"
sidebar_current: "docs-nexus-datasource-repository"
description: |-
  Use this data source to get a repository data structure
---

# nexus_repository

Use this data source to get a repository data structure

## Example Usage

```hcl
data "nexus_repository" "maven-central" {
  name = "maven-central"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique identifier for this repository
* `apt_signing` - (Optional) Apt signing configuration for the repository
* `apt` - (Optional) Apt specific configuration of the repository
* `cleanup` - (Optional) Cleanup policies
* `docker` - (Optional) Docker specific configuration of the repository
* `format` - (Optional) Repository format
* `group` - (Optional) Configuration for repository group
* `http_client` - (Optional) HTTP Client configuration for proxy repositories
* `maven` - (Optional) Maven specific configuration of the repository
* `negative_cache` - (Optional) Configuration of the negative cache handling
* `online` - (Optional) Whether this repository accepts incoming requests
* `proxy` - (Optional) Configuration for the proxy repository
* `storage` - (Optional) The storage configuration of the repository
* `type` - (Optional) Repository type

The `apt` object supports the following:

* `distribution` - (Required) The linux distribution

The `apt_signing` object supports the following:

* `keypair` - (Required) PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor )
* `passphrase` - (Required) Passphrase for the keypair

The `cleanup` object supports the following:

* `policy_names` - (Required) List of policy names

The `docker` object supports the following:

* `force_basic_auth` - (Optional) Whether to force authentication (Docker Bearer Token Realm required if false)
* `http_port` - (Optional) Create an HTTP connector at specified port
* `https_port` - (Optional) Create an HTTPS connector at specified port
* `v1enabled` - (Optional) Whether to allow clients to use the V1 API to interact with this repository

The `group` object supports the following:

* `member_names` - (Required) Member repositories names

The `http_client` object supports the following:

* `authentication` - (Optional) Authentication configuration of the HTTP client
* `auto_block` - (Optional) Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive
* `blocked` - (Optional) Whether to block outbound connections on the repository
* `connection` - (Optional) Connection configuration of the HTTP client

The `authentication` object supports the following:

* `type` - (Required) Authentication type
* `ntlm_domain` - (Optional) The ntlm domain to connect
* `ntlm_host` - (Optional) The ntlm host to connect
* `username` - (Optional) The username used by the proxy repository

The `connection` object supports the following:

* `retries` - (Optional) Total retries if the initial connection attempt suffers a timeout
* `timeout` - (Optional) Seconds to wait for activity before stopping and retrying the connection

The `maven` object supports the following:

* `layout_policy` - (Optional) Validate that all paths are maven artifact or metadata paths
* `version_policy` - (Optional) What type of artifacts does this repository store

The `negative_cache` object supports the following:

* `enabled` - (Optional) Whether to cache responses for content not present in the proxied repository
* `ttl` - (Optional) How long to cache the fact that a file was not found in the repository (in minutes)

The `proxy` object supports the following:

* `remote_url` - (Required) Location of the remote repository being proxied
* `content_max_age` - (Optional) How long (in minutes) to cache artifacts before rechecking the remote repository
* `metadata_max_age` - (Optional) How long (in minutes) to cache metadata before rechecking the remote repository.

The `storage` object supports the following:

* `blob_store_name` - (Optional) Blob store used to store repository contents
* `strict_content_type_validation` - (Optional) Whether to validate uploaded content's MIME type appropriate for the repository format
* `write_policy` - (Optional) Controls if deployments of and updates to assets are allowed


