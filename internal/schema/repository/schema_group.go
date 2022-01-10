package repository

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceGroup = &schema.Schema{
		Description: "Configuration for repository group",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"member_names": {
					Description: "Member repositories names",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Required: true,
					Set: func(v interface{}) int {
						return schema.HashString(strings.ToLower(v.(string)))
					},
					Type: schema.TypeSet,
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeList,
	}
)
