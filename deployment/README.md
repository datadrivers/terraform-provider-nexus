# Provider Deployment

This folder has the provider deployment bits that are needed.

You'll need to follow [this link](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers#publishing-a-provider) for docs on what to do.


## Steps

1. Run `create-provider.sh`
2. Bump the version in version.json
3. Run `create-provider-version.sh`
4. Run `get-sigs.sh`
5. Run `curl -T <file> <url>` for URLs returned by `create-provider-version.sh` to upload signatures
6. Check `platform.json` and `create-platform.sh` for version numbers
7. Run `create-platform.sh`
8. Run `curl -T <file> <url>` with returned URL to upload the binary.
9. Navigate to https://app.terraform.io/app/octopus/registry/providers/private/octopus/nexus/latest/overview and go to the version. If all is well, you should not see any warning banners. If there's a banner, click Manage Provider -> Show Release Files
