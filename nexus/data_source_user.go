/*
This data source is deprecated. Please use the data source "nexus_security_user" instead.

Use this data source to get a user data structure

Example Usage

```hcl
data "nexus_user" "admin" {
  userid = "admin"
}
```
*/
package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
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

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("userid").(string))

	return resourceUserRead(d, m)
}
