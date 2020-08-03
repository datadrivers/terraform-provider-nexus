package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecurityLDAPOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityLDAPOrderCreate,
		Read:   resourceSecurityLDAPOrderRead,
		Update: resourceSecurityLDAPOrderUpdate,
		Delete: resourceSecurityLDAPOrderDelete,

		Schema: map[string]*schema.Schema{
			"order": {
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
	client := m.(nexus.Client)
	order := interfaceSliceToStringSlice(d.Get("order").([]interface{}))
	if err := client.LDAPChangeOrder(order); err != nil {
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
