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
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAnonymous() *schema.Resource {
	return &schema.Resource{
		Create: resourceAnonymousUpdate,
		Read:   resourceAnonymousRead,
		Update: resourceAnonymousUpdate,
		Delete: resourceAnonymousDelete,
		Exists: resourceAnonymousExists,
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

func getAnonymousFromResourceData(d *schema.ResourceData) nexus.AnonymousConfig {
	return nexus.AnonymousConfig{
		Enabled:   d.Get("enabled").(bool),
		UserID:    d.Get("user_id").(string),
		RealmName: d.Get("realm_name").(string),
	}
}

func setAnonymousToResourceData(anonymous *nexus.AnonymousConfig, d *schema.ResourceData) error {
	d.SetId("anonymous")
	d.Set("enabled", anonymous.Enabled)
	d.Set("user_id", anonymous.UserID)
	d.Set("realm_name", anonymous.RealmName)
	return nil
}

func resourceAnonymousRead(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	anonymous, err := client.AnonymousRead()
	if err != nil {
		return err
	}

	return setAnonymousToResourceData(anonymous, d)
}

func resourceAnonymousUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	anonymous := getAnonymousFromResourceData(d)
	if err := client.AnonymousUpdate(anonymous); err != nil {
		return err
	}

	return resourceAnonymousRead(d, m)
}

func resourceAnonymousDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAnonymousExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(nexus.Client)

	privilege, err := client.AnonymousRead()
	return privilege != nil, err
}
