package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a user data structure.",

		Read: dataSourceSecurityUserRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"userid": {
				Description: "The userid which is required for login",
				Type:        schema.TypeString,
				Required:    true,
			},
			"firstname": {
				Description: "The first name of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"lastname": {
				Description: "The last name of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The email address associated with the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"roles": {
				Description: "The roles which the user has been assigned within Nexus.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Description: "The user's status, e.g. active or disabled.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSecurityUserRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("userid").(string))

	return resourceSecurityUserRead(d, m)
}
