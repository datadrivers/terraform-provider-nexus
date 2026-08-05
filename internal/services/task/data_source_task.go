package task

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTaskRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaskRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("id").(string))
	d.SetId(d.Get("name").(string))
	d.SetId(d.Get("type").(string))
	d.SetId(d.Get("current_state").(string))
	d.SetId(d.Get("last_run_result").(string))
	d.SetId(d.Get("next_run").(string))
	return resourceTaskRead(d, m)
}
