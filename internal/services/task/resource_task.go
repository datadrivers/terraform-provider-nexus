package task

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema/task"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTaskCreate,
		Read:   resourceTaskRead,
		Update: resourceTaskUpdate,
		Delete: resourceTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"alert_email": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"notification_condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frequency": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_date": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"time_zone_offset": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"recurring_days": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTaskCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	frequencySet := d.Get("frequency").(*schema.Set)
	frequencyMap := frequencySet.List()[0].(map[string]interface{})
	task := nexusSchema.TaskCreateStruct{
		Name:                  d.Get("name").(string),
		Type:                  d.Get("type").(string),
		Enabled:               d.Get("enabled").(bool),
		AlertEmail:            d.Get("alert_email").(string),
		NotificationCondition: d.Get("notification_condition").(string),
		Frequency: nexusSchema.FrequencyXO{
			Schedule:       frequencyMap["schedule"].(string),
			StartDate:      frequencyMap["start_date"].(int),
			TimeZoneOffset: frequencyMap["time_zone_offset"].(string),
			RecurringDays:  frequencyMap["recurring_days"].([]interface{}),
			CronExpression: frequencyMap["cron_expression"].(string),
		},
		Properties: d.Get("properties").(map[string]interface{}),
	}

	newTask, err := client.Task.CreateTask(&task)
	if err != nil {
		return err
	}
	d.SetId(newTask.ID)
	return resourceTaskRead(d, m)
}

func resourceTaskRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	id := d.Id()

	task, err := client.Task.GetTask(id)
	if err != nil {
		return err
	}

	d.Set("name", task.Name)
	d.Set("type", task.Type)

	return nil
}

func resourceTaskUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	if d.HasChange("name") || d.HasChange("enabled") || d.HasChange("Frequency") || d.HasChange("notification_condition") {
		task := getTaskFromResourceData(d)
		if err := client.Task.UpdateTask(d.Id(), &task); err != nil {
			return err
		}
	}
	task := getTaskFromResourceData(d)

	err := client.Task.UpdateTask(d.Id(), &task)
	if err != nil {
		return err
	}
	return resourceTaskRead(d, m)
}

func getTaskFromResourceData(d *schema.ResourceData) nexusSchema.TaskCreateStruct {
	frequencySet := d.Get("frequency").(*schema.Set)
	frequencyMap := frequencySet.List()[0].(map[string]interface{})
	return nexusSchema.TaskCreateStruct{
		Name:                  d.Get("name").(string),
		Enabled:               d.Get("enabled").(bool),
		AlertEmail:            d.Get("alert_email").(string),
		NotificationCondition: d.Get("notification_condition").(string),
		Frequency: nexusSchema.FrequencyXO{
			Schedule:       frequencyMap["schedule"].(string),
			StartDate:      frequencyMap["start_date"].(int),
			TimeZoneOffset: frequencyMap["time_zone_offset"].(string),
			RecurringDays:  frequencyMap["recurring_days"].([]interface{}),
			CronExpression: frequencyMap["cron_expression"].(string),
		},
	}
}

func resourceTaskDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Task.DeleteTask(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
