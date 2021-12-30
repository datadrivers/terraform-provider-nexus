/*
Use this resource to change the LDAP order.

Example Usage

Set order of LDAP server

```hcl
resource "nexus_security_ldap" "server1" {
  ...
  name = "server1"
}

resource "nexus_security_ldap" "server2" {
  ...
  name = "server2"
}

resource "nexus_security_ldap_order" {
  order = [
    nexus_security_ldap.server1.name,
    nexus_security_ldap.server2.name,
  ]
}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityLDAPOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityLDAPOrderCreate,
		Read:   resourceSecurityLDAPOrderRead,
		Update: resourceSecurityLDAPOrderUpdate,
		Delete: resourceSecurityLDAPOrderDelete,

		Schema: map[string]*schema.Schema{
			"order": {
				Description: "Ordered list of LDAP server",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				Type:     schema.TypeList,
			},
		},
	}
}

func resourceSecurityLDAPOrderCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	order := interfaceSliceToStringSlice(d.Get("order").([]interface{}))
	if err := client.Security.LDAP.ChangeOrder(order); err != nil {
		return err
	}

	d.SetId("change-order")
	d.Set("order", order)

	return resourceSecurityLDAPOrderRead(d, m)
}

func resourceSecurityLDAPOrderRead(d *schema.ResourceData, m interface{}) error {
	// d.Set("order", d.Get("order"))

	return nil
}

func resourceSecurityLDAPOrderUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSecurityLDAPOrderCreate(d, m)
}

func resourceSecurityLDAPOrderDelete(d *schema.ResourceData, m interface{}) error {
	// return fmt.Errorf("Nexus API does not support deleting LDAP server order")
	return nil
}
