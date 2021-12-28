/*
Use this resource to activate Nexus Security LDAP and order

Example Usage

```hcl
resource "nexus_security_realms" "example" {
  active = [
	"NexusAuthenticatingRealm",
	"NexusAuthorizingRealm",
  ]
}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecurityRealms() *schema.Resource {
	return &schema.Resource{
		Create: resourceRealmsCreate,
		Read:   resourceRealmsRead,
		Update: resourceRealmsUpdate,
		Delete: resourceRealmsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Description: "The realm IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				Type:     schema.TypeList,
			},
		},
	}
}

func resourceRealmsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	realmIDs := interfaceSliceToStringSlice(d.Get("active").([]interface{}))
	if err := client.Security.Realm.Activate(realmIDs); err != nil {
		return err
	}

	return resourceRealmsRead(d, m)
}

func resourceRealmsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	activeRealms, err := client.Security.Realm.ListActive()
	if err != nil {
		return err
	}

	d.SetId("active")
	if err := d.Set("active", stringSliceToInterfaceSlice(activeRealms)); err != nil {
		return err
	}

	return nil
}

func resourceRealmsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceRealmsCreate(d, m)
}

func resourceRealmsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
