package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRealmsActive() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRealmsActiveRead,

		Schema: map[string]*schema.Schema{
			"realms": {
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceRealmsActiveRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	activeRealms, err := nexusClient.RealmsActive()
	if err != nil {
		return err
	}

	d.SetId("active")
	if err := d.Set("realms", stringSliceToInterfaceSlice(activeRealms)); err != nil {
		return err
	}

	return nil
}
