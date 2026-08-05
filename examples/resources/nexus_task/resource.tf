resource "nexus_task" "default" {
    name = "Cleanup service task"
    type = "tags.cleanup"
    enabled = true
    alert_email = "test@test.com"
    notification_condition = "FAILURE"
    frequency {
        schedule = "daily"
        cron_expression = null
        recurring_days = null
        start_date = 1740122615
    }
    properties = {
    }
}
