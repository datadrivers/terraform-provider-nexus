/*
Use this data source to list all security realms.

Example Usage

```hcl
data "nexus_security_realms" "default" {}
```
*/
package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityRealms() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to list all security realms.",

		Read: dataSourceRealmsRead,
		Schema: map[string]*schema.Schema{
			"active": {
				Description: "List of active security realms",
				Computed:    true,
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
				Description: "List of available security realms",
				Computed:    true,
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
	client := m.(*nexus.NexusClient)

	availableRealms, err := client.Security.Realm.ListAvailable()
	if err != nil {
		return err
	}

	activeRealmIDs, err := client.Security.Realm.ListActive()
	if err != nil {
		return err
	}

	activeRealms := make([]security.Realm, len(activeRealmIDs))
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

func flattenSecurityRealms(realms []security.Realm) []map[string]interface{} {
	data := make([]map[string]interface{}, len(realms))
	for k, v := range realms {
		data[k] = map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
		}
	}

	return data
}
