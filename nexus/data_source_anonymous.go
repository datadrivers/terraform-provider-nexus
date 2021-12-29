/*
Use this get the anonymous configuration of the nexus repository manager.

Example Usage

```hcl
data "nexus_anonymous" "nexus" {
}

output "nexus_anonymous_enabled" {
  description = "Anonymous enabled?"
  value       = data.nexus_anonymous.nexus.enabled
}
```
*/
package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAnonymous() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAnonymousRead,
		Schema: map[string]*schema.Schema{
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
