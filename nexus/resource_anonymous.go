/*
Use this resource to change the anonymous configuration of the nexus repository manager.

Example Usage

```hcl
resource "nexus_anonymous" "example" {
  enabled = true
  user_id = "exampleUser"
}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAnonymous() *schema.Resource {
	return &schema.Resource{
		Create: resourceAnonymousUpdate,
		Read:   resourceAnonymousRead,
		Update: resourceAnonymousUpdate,
		Delete: resourceAnonymousDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Activate the anonymous access to the repository manager. Default: false",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"user_id": {
				Description: "The user id used by anonymous access. Default: \"anonymous\"",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "anonymous",
			},
			"realm_name": {
				Description: "The name of the used realm. Default: \"NexusAuthorizingRealm\"",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NexusAuthorizingRealm",
			},
		},
	}
}

func getAnonymousFromResourceData(d *schema.ResourceData) security.AnonymousAccessSettings {
	return security.AnonymousAccessSettings{
		Enabled:   d.Get("enabled").(bool),
		UserID:    d.Get("user_id").(string),
		RealmName: d.Get("realm_name").(string),
	}
}

func setAnonymousToResourceData(anonymous *security.AnonymousAccessSettings, d *schema.ResourceData) error {
	d.SetId("anonymous")
	d.Set("enabled", anonymous.Enabled)
	d.Set("user_id", anonymous.UserID)
	d.Set("realm_name", anonymous.RealmName)
	return nil
}

func resourceAnonymousRead(d *schema.ResourceData, m interface{}) error {
	service := m.(nexus.NexusService)

	anonymous, err := service.Security.Anonymous.Read()
	if err != nil {
		return err
	}

	return setAnonymousToResourceData(anonymous, d)
}

func resourceAnonymousUpdate(d *schema.ResourceData, m interface{}) error {
	service := m.(nexus.NexusService)

	anonymous := getAnonymousFromResourceData(d)
	if err := service.Security.Anonymous.Update(anonymous); err != nil {
		return err
	}

	return resourceAnonymousRead(d, m)
}

func resourceAnonymousDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
