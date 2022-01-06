package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	DataSourceID = &schema.Schema{
		Description: "Used to identify data source at nexus",
		Computed:    true,
		Type:        schema.TypeString,
	}
	ResourceID = &schema.Schema{
		Description: "Used to identify resource at nexus",
		Computed:    true,
		Type:        schema.TypeString,
	}
)
