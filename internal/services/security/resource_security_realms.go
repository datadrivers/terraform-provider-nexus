package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityRealms() *schema.Resource {
	return &schema.Resource{
		Description: `Use this resource to activate and order the Nexus Security realms.

!> This resource can only be used **once** for a nexus`,

		Create: resourceRealmsCreate,
		Read:   resourceRealmsRead,
		Update: resourceRealmsUpdate,
		Delete: resourceRealmsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"active": {
				Description: "Set the active security realms in the order they should be used.",
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
	realmIDs := tools.InterfaceSliceToStringSlice(d.Get("active").([]interface{}))
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
	if err := d.Set("active", tools.StringSliceToInterfaceSlice(activeRealms)); err != nil {
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
