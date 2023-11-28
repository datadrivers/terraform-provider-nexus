package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ResourceDocker = &schema.Schema{
		Description: "docker contains the configuration of the docker repository",
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"force_basic_auth": {
					Description: "Whether to force authentication (Docker Bearer Token Realm required if false)",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"http_port": {
					Description: "Create an HTTP connector at specified port",
					Optional:    true,
					Type:        schema.TypeInt,
				},
				"https_port": {
					Description: "Create an HTTPS connector at specified port",
					Optional:    true,
					Type:        schema.TypeInt,
				},
				"v1_enabled": {
					Description: "Whether to allow clients to use the V1 API to interact with this repository",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"subdomain": {
					Description: "Use sub-domain routing for this repository",
					Required:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}
	DataSourceDocker = &schema.Schema{
		Description: "docker contains the configuration of the docker repository",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"force_basic_auth": {
					Description: "Whether to force authentication (Docker Bearer Token Realm required if false)",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"http_port": {
					Description: "Create an HTTP connector at specified port",
					Computed:    true,
					Type:        schema.TypeInt,
				},
				"https_port": {
					Description: "Create an HTTPS connector at specified port",
					Computed:    true,
					Type:        schema.TypeInt,
				},
				"v1_enabled": {
					Description: "Whether to allow clients to use the V1 API to interact with this repository",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"subdomain": {
					Description: "Use sub-domain routing for this repository",
					Required:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}
)
