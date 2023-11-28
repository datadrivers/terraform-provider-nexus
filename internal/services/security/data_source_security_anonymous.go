package security

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityAnonymous() *schema.Resource {
	return &schema.Resource{
		Description: "Use this to get the anonymous configuration of the nexus repository manager.",

		Read: dataSourceSecurityAnonymousRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"enabled": {
				Computed:    true,
				Description: "Activate the anonymous access to the repository manager",
				Type:        schema.TypeBool,
			},
			"user_id": {
				Computed:    true,
				Description: "The user id used by anonymous access",
				Type:        schema.TypeString,
			},
			"realm_name": {
				Computed:    true,
				Description: "The name of the used realm",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceSecurityAnonymousRead(d *schema.ResourceData, m interface{}) error {
	return resourceSecurityAnonymousRead(d, m)
}
