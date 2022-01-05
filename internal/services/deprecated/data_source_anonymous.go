package deprecated

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAnonymous() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source is deprecated. Please use the data source nexus_security_anonymous instead.",
		Description: `!> This data source is deprecated. Please use the data source "nexus_security_anonymous" instead.

Use this get the anonymous configuration of the nexus repository manager.`,

		Read: dataSourceAnonymousRead,
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

func dataSourceAnonymousRead(d *schema.ResourceData, m interface{}) error {
	return resourceAnonymousRead(d, m)
}
