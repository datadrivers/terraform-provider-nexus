/*
Use this resource to create a Nexus Security LDAP

Example Usage

```hcl
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
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
	nexusClient := m.(nexus.Client)
	realmIDs := interfaceSliceToStringSlice(d.Get("active").([]interface{}))
	if err := nexusClient.RealmsActivate(realmIDs); err != nil {
		return err
	}

	return resourceRealmsRead(d, m)
}

func resourceRealmsRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	activeRealms, err := nexusClient.RealmsActive()
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
