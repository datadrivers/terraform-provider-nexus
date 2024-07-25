#!/bin/bash

# Omit the leading 'v' if present:
version="2.4.0-rc2-bump"

wget "https://github.com/octoenergy/terraform-provider-nexus/releases/download/v${version}/terraform-provider-nexus_${version}_SHA256SUMS"
wget "https://github.com/octoenergy/terraform-provider-nexus/releases/download/v${version}/terraform-provider-nexus_${version}_SHA256SUMS.sig"
wget "https://github.com/octoenergy/terraform-provider-nexus/releases/download/v${version}/terraform-provider-nexus_${version}_linux_amd64.zip"
