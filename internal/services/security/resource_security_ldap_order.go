package security

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/dre2004/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityLDAPOrder() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to change the LDAP order.",

		Create: resourceSecurityLDAPOrderCreate,
		Read:   resourceSecurityLDAPOrderRead,
		Update: resourceSecurityLDAPOrderUpdate,
		Delete: resourceSecurityLDAPOrderDelete,

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
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
	order := tools.InterfaceSliceToStringSlice(d.Get("order").([]interface{}))
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
