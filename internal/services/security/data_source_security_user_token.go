/*
Use this data source to get the global user-token configuration.

---
**PRO Feature**
---

Example Usage

```hcl
data "nexus_security_user_token" "nexus" {}

output "nexus_user_token_enabled" {
  description = "User Tokens enabled?"
  value       = data.nexus_security_user_token.nexus.enabled
}
```
*/
package security

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityUserToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityUserTokenRead,

		Schema: map[string]*schema.Schema{
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
