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
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTaskRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("id").(string))
	return resourceTaskRead(d, m)
}
