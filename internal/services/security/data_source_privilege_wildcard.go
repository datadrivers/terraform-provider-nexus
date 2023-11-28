package security

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityPrivilegeWildcard() *schema.Resource {
	return &schema.Resource{
		Description: `Use this data source to get a privilege for an wildcard`,

		Read: dataSourceSecurityPrivilegeWildcardRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"name": {
				Description: "Name of the wildcard privilege",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the wildcard privilege",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"readonly": {
				Description: "Whether it is readonly or not",
				Computed:    true,
				Type:        schema.TypeBool,
			},
			"pattern": {
				Description: "Pattern of the wildcard privilege",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceSecurityPrivilegeWildcardRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeWildcardRead(d, m)
}
