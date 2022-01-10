package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceProxy = &schema.Schema{
		Description: "Configuration for the proxy repository",
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"content_max_age": {
					Description: "How long (in minutes) to cache artifacts before rechecking the remote repository",
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     1440,
				},
				"metadata_max_age": {
					Description: "How long (in minutes) to cache metadata before rechecking the remote repository.",
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     1440,
				},
				"remote_url": {
					Description: "Location of the remote repository being proxied",
					Type:        schema.TypeString,
					Optional:    true,
				},
			},
		},
	}
)
