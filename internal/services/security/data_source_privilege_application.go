package security

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityPrivilegeApplication() *schema.Resource {
	return &schema.Resource{
		Description: `Use this data source to get a privilege for an application`,

		Read: dataSourceSecurityPrivilegeApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"name": {
				Description: "Name of the application privilege",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the application privilege",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"readonly": {
				Description: "Whether it is readonly or not",
				Computed:    true,
				Type:        schema.TypeBool,
			},
			"actions": {
				Description: "Description of the application privilege",
				Computed:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"domain": {
				Description: "Domain of the application privilege",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceSecurityPrivilegeApplicationRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeApplicationRead(d, m)
}
