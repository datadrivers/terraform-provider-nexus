package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceComponent = &schema.Schema{
		DefaultFunc: repositoryComponentDefault,
		Description: "Component configuration for the hosted repository",
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"proprietary_components": {
					Description: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
					Type:        schema.TypeBool,
					Required:    true,
				},
			},
		},
	}

	DataSourceComponent = &schema.Schema{
		Description: "Component configuration for the hosted repository",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"proprietary_components": {
					Description: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
					Type:        schema.TypeBool,
					Computed:    true,
				},
			},
		},
	}
)

func repositoryComponentDefault() (interface{}, error) {
	return []map[string]interface{}{
		{
			"proprietary_components": false,
		},
	}, nil
}
