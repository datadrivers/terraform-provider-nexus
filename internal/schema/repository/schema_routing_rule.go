package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceRoutingRule = &schema.Schema{
		Description: "The name of the routing rule assigned to this repository",
		Optional:    true,
		Type:        schema.TypeString,
	}
	DataSourceRoutingRule = &schema.Schema{
		Description: "The name of the routing rule assigned to this repository",
		Computed:    true,
		Type:        schema.TypeString,
	}
)
