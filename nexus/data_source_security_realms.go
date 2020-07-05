package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSecurityRealms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRealmsRead,

		Schema: map[string]*schema.Schema{
			"active": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed:    true,
							Description: "Realm ID",
							Type:        schema.TypeString,
						},
						"name": {
							Computed:    true,
							Description: "Realm name",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			"available": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed:    true,
							Description: "Realm ID",
							Type:        schema.TypeString,
						},
						"name": {
							Computed:    true,
							Description: "Realm name",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceRealmsRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	availableRealms, err := nexusClient.RealmsAvailable()
	if err != nil {
		return err
	}

	activeRealmIDs, err := nexusClient.RealmsActive()
	if err != nil {
		return err
	}

	activeRealms := make([]nexus.Realm, len(activeRealmIDs))
	for i, activeRealmID := range activeRealmIDs {
		for _, v := range availableRealms {
			if v.ID == activeRealmID {
				activeRealms[i] = v
				break
			}
		}
	}

	d.SetId("security-realms")
	if err := d.Set("active", flattenSecurityRealms(activeRealms)); err != nil {
		return err
	}
	if err := d.Set("available", flattenSecurityRealms(availableRealms)); err != nil {
		return err
	}

	return nil
}

func flattenSecurityRealms(realms []nexus.Realm) []map[string]interface{} {
	data := make([]map[string]interface{}, len(realms))
	for k, v := range realms {
		data[k] = map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
		}
	}

	return data
}
