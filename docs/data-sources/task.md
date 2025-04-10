---
page_title: "Data Source nexus_task"
subcategory: "Task"
description: |-
  Use this data source to get a task data structure.
---
# Data Source nexus_task
Use this data source to get a task data structure.

## Example Usage
```terraform
data "nexus_task" "example" {
  id = "example-task-id"
}
```

## Schema

### Required
- `id` (String) The ID of the task.

### Read-Only

- `name` (String) The name of the task.
- `type` (String) The type of the task.
- `current_state` (String) The current state of the task.