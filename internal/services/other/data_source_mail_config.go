package other

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEmailConfig() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to work with the e-mail config",

		Read: DataSourceEmailConfigRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"enabled": {
				Description: "Whether the mail config is active or not",
				Computed:    true,
				Type:        schema.TypeBool,
			},
			"host": {
				Description: "Host",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func DataSourceEmailConfigRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("id").(string))
	return resourceEmailConfigRead(d, m)
}
