package repository

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceCleanup = &schema.Schema{
		Description: "Cleanup policies",
		Type:        schema.TypeList,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_names": {
					Description: "List of policy names",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
					Set: func(v interface{}) int {
						return schema.HashString(strings.ToLower(v.(string)))
					},
					Type: schema.TypeSet,
				},
			},
		},
	}
	DataSourceCleanup = &schema.Schema{
		Description: "Cleanup policies",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_names": {
					Description: "List of policy names",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Computed: true,
					Set: func(v interface{}) int {
						return schema.HashString(strings.ToLower(v.(string)))
					},
					Type: schema.TypeSet,
				},
			},
		},
	}
)
