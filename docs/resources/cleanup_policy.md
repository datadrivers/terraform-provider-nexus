---
page_title: "Resource nexus_cleanup_policy"
subcategory: "cleanup_policy"
description: |-
  Use this resource to manage users.
---
# Resource nexus_cleanup_policy
Use this resource to manage cleanup policies.

## Example Usage
```terraform
resource "nexus_cleanup_policy" "example" {
  name                      = "maven-artifacts-clean-up"
  notes                     = "This is an example cleanup policy"
  criteria_last_downloaded  = 0
  criteria_release_type     = "RELEASES"
  criteria_asset_regex      = ".*"
  retain                    = null
  format                    = "maven2"
}
```

## Schema

### Required

- `name` (String) The name of the cleanup policy.
- `format` (String) The format of the repository.

### Optional

- `notes` (String) Notes for the cleanup policy.
- `criteria_last_blob_updated` (Int) The criteria for the last blob updated.
- `criteria_last_downloaded` (Int) The criteria for the last downloaded.
- `criteria_release_type` (String) The criteria for the release type.
- `criteria_asset_regex` (String) The criteria for the asset regex.
- `retain` (Int) The number of items to retain.

### Read-Only

- `id` (String) The ID of the cleanup policy.

## Import

Import is supported using the following syntax:
```shell
# import using the name of the cleanup policy
terraform import nexus_cleanup_policy.example <name_cleanup_policy>
```
