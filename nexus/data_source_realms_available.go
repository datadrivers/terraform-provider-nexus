package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRealmsAvailable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRealmsAvailableRead,

		Schema: map[string]*schema.Schema{
			"realms": {
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

func dataSourceRealmsAvailableRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	realms, err := nexusClient.RealmsAvailable()
	if err != nil {
		return err
	}

	d.SetId("available")
	if err := d.Set("realms", flattenAvailableRealms(realms)); err != nil {
		return err
	}

	return nil
}

func flattenAvailableRealms(realms []nexus.Realm) []map[string]interface{} {
	data := make([]map[string]interface{}, len(realms))
	for k, v := range realms {
		data[k] = map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
		}
	}

	return data
}
