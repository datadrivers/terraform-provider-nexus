package other

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRoutingRule() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to work with routing rules.",

		Read: dataSourceRoutingRuleRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"name": {
				Description: "The name of the routing rule",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Computed:    true,
				Description: "The description of the routing rule",
				Type:        schema.TypeString,
			},
			"mode": {
				Computed:    true,
				Description: "The mode describe how to hande with mathing requests. Possible values: `BLOCK` or `ALLOW`",
				Type:        schema.TypeString,
			},
			"matchers": {
				Computed:    true,
				Description: "Matchers is a list of regular expressions used to identify request paths that are allowed or blocked (depending on above mode)",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
		},
	}
}

func dataSourceRoutingRuleRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceRoutingRuleRead(d, m)
}
