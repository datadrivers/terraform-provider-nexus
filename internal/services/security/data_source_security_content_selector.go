package security

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityContentSelector() *schema.Resource {
	return &schema.Resource{
		Description: "Use this to get a specified content selector.",

		Read: dataSourceSecurityContentSelectorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Content selector name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A description of the content selector",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"expression": {
				Description: "The content selector expression",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceSecurityContentSelectorRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceSecurityContentSelectorRead(d, m)
}
