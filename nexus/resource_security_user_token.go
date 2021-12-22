/*
Use this resource to manage the global configuration for the user-tokens

---
**PRO Feature**
---

Example Usage

```hcl
resource "nexus_security_user_token" "nexus" {
    enabled         = true
	protect_content = false
}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecurityUserToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityUserTokenUpdate,
		Read:   resourceSecurityUserTokenRead,
		Update: resourceSecurityUserTokenUpdate,
		Delete: resourceSecurityUserTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Activate the feature of user tokens.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"protect_content": {
				Description: "Require user tokens for repository authentication. This does not effect UI access.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func getSecurityUserTokenFromResourceData(d *schema.ResourceData) nexus.UserTokenConfiguration {
	return nexus.UserTokenConfiguration{
		Enabled:        d.Get("enabled").(bool),
		ProtectContent: d.Get("protect_content").(bool),
	}
}

func setSecurityUserTokenToResourceData(token *nexus.UserTokenConfiguration, d *schema.ResourceData) {
	d.SetId("golbalUserTokenConfiguration")
	d.Set("enabled", token.Enabled)
	d.Set("protect_content", token.ProtectContent)
}

func resourceSecurityUserTokenRead(d *schema.ResourceData, m interface{}) error {
	service := m.(nexus.NexusService)
	token, err := nexusClient.UserTokensRead()
	if err != nil {
		return err
	}
	setSecurityUserTokenToResourceData(token, d)
	return nil
}

func resourceSecurityUserTokenUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(nexus.NexusService)

	token := getSecurityUserTokenFromResourceData(d)
	if err := nexusClient.UserTokensApply(token); err != nil {
		return err
	}

	return resourceSecurityUserTokenRead(d, m)
}

func resourceSecurityUserTokenDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
