package security

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityUserToken() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this data source to get the global user-token configuration.`,

		Read: dataSourceSecurityUserTokenRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"enabled": {
				Computed:    true,
				Description: "Activate the feature of user tokens.",
				Type:        schema.TypeBool,
			},
			"protect_content": {
				Computed:    true,
				Description: "Require user tokens for repository authentication. This does not effect UI access.",
				Type:        schema.TypeBool,
			},
		},
	}
}

func dataSourceSecurityUserTokenRead(d *schema.ResourceData, m interface{}) error {
	return resourceSecurityUserTokenRead(d, m)
}
