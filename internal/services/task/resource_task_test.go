package task_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/task"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTask_basic(t *testing.T) {
	resourceName := "nexus_task.test"

	task := task.TaskCreateStruct{
		Name:                  acctest.RandString(20),
		Type:                  "tags.cleanup",
		Enabled:               true,
		AlertEmail:            acctest.RandString(20) + "@example.com",
		NotificationCondition: "FAILURE",
		Frequency: task.FrequencyXO{
			StartDate:      1740122615,
			Schedule:       "daily",
			TimeZoneOffset: "+1",
			//RecurringDays:  []{1, 2, 3},
			CronExpression: "0 0 * * *",
		},
		Properties: map[string]interface{}{},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTaskCreateConfig(task),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", task.Name),
					resource.TestCheckResourceAttr(resourceName, "type", task.Type),
					resource.TestCheckResourceAttr(resourceName, "alert_email", task.AlertEmail),
					resource.TestCheckResourceAttr(resourceName, "notification_condition", task.NotificationCondition),
					resource.TestCheckResourceAttr(resourceName, "frequency.0.schedule", task.Frequency.Schedule),
					resource.TestCheckResourceAttr(resourceName, "frequency.0.start_date", strconv.FormatInt(int64(task.Frequency.StartDate), 10)),
					resource.TestCheckResourceAttr(resourceName, "frequency.0.time_zone_offset", task.Frequency.TimeZoneOffset),
					resource.TestCheckResourceAttr(resourceName, "frequency.0.cron_expression", task.Frequency.CronExpression),
				),
			},
		},
	})
}

func testAccResourceTaskCreateConfig(task task.TaskCreateStruct) string {
	return fmt.Sprintf(`
resource "nexus_task" "test" {
	name = "%s"
	type = "%s"
	enabled = %t
	alert_email = "%s"
	notification_condition = "%s"
	frequency {
		schedule = "%s"
		start_date = %d
		time_zone_offset = "%s"
		recurring_days = [%s]
		cron_expression = "%s"
	}
	properties = {}
}
`, task.Name, task.Type, task.Enabled, task.AlertEmail, task.NotificationCondition, task.Frequency.Schedule, task.Frequency.StartDate, task.Frequency.TimeZoneOffset, formatRecurringDays(task.Frequency.RecurringDays), task.Frequency.CronExpression)
}

func formatRecurringDays(days []interface{}) string {
	var formattedDays []string
	for _, day := range days {
		formattedDays = append(formattedDays, strconv.Itoa(day.(int)))
	}
	return strings.Join(formattedDays, ", ")
}
