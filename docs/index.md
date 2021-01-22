---
layout: "nexus"
page_title: "Provider: Nexus"
sidebar_current: "docs-nexus-index"
description: |-
  The Nexus provider allows Terraform to read from, write to, and configure Sonatype Nexus Repository Manager
---

# Nexus Provider
  
The Nexus provider allows Terraform to read from, write to, and configure [Sonatype Nexus Repository Manager](https://www.sonatype.com/product-nexus-repository).

-> **Note** This provider hat been implemented and tested with Sonatype Nexus Repository Manager OSS `3.29.0-02`.

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

## Arguments Reference

The following arguments are supported:

* `insecure`  - (Optional) Boolean to specify wether insecure SSL connections are allowed or not.
* `password`  - (Required) Password of user to connect to API.
* `url`       - (Required) URL of Nexus to reach API.
* `usernamme` - (Required) Username used to connect to API.

## Author

[Datadrivers GmbH](https://www.datadrivers.de)
