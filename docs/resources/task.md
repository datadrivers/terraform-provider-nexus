---
page_title: "Resource nexus_task"
subcategory: "Task"
description: |-
  Use this resource to manage tasks.
---
# Resource nexus_task
Use this resource to manage tasks.

## Example Usage
```terraform
resource "nexus_task" "example" {
  name                   = "example-task"
  type                   = "repository.cleanup"
  enabled                = true
  alert_email            = "alert@example.com"
  notification_condition = "FAILURE"
  frequency {
    schedule       = "daily"
    start_date     = 1633046400000
    time_zone_offset = "UTC+0"
    recurring_days = [1, 2, 3]
    cron_expression = "0 0 * * *"
  }
  properties = {}
}
```

## Schema

### Required

- `name` (String) The name of the task.
- `type` (String) The type of the task.
- `enabled` (Boolean) Whether the task is enabled or not.
- `alert_email` (String) The email address to send alerts to.
- `notification_condition` (String) The condition under which to send notifications.
- `frequency` (Block) The frequency of the task.
  - `schedule` (String) The schedule of the task.
  - `start_date` (Number) The start date of the task.
  - `time_zone_offset` (String) The time zone offset of the task.
  - `recurring_days` (List of Number) The recurring days of the task.
  - `cron_expression` (String) The cron expression of the task.

### Optional

- `properties` (Map of String) The properties of the task.

### Read-Only

- `id` (String) The ID of the task.


### Frequency Schema
- `schedule` (String) The schedule for the task.
- `start_date` (Int) The start date for the task.
- `time_zone_offset` (String) The time zone offset for the task.
- `recurring_days` (List of Int) The recurring days for the task.
- `cron_expression` (String) The cron expression for the task.

## Import

Import is supported using the following syntax:
```shell
# import using the ID of the task
terraform import nexus_task.example <task_id>
```

