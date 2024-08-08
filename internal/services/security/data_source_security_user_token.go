package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
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
				Description: "Activation of the user tokens feature.",
				Type:        schema.TypeBool,
			},
			"protect_content": {
				Computed:    true,
				Description: "Require user tokens for repository authentication. This does not effect UI access.",
				Type:        schema.TypeBool,
			},
			"expiration_enabled": {
				Description: "Activation of the user tokens expiration feature.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"expiration_days": {
				Description: "Number of days user tokens remain valid.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceSecurityUserTokenRead(d *schema.ResourceData, m interface{}) error {
	return resourceSecurityUserTokenRead(d, m)
}
