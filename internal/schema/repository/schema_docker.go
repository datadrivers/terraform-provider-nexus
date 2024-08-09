package repository

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
					Description: "Pro-only: Whether to allow clients to use subdomain routing connector",
					Optional:    true,
					Type:        schema.TypeString,
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
					Description: "Pro-only: Whether to allow clients to use subdomain routing connector",
					Computed:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	ResourceDockerHostedStorage = &schema.Schema{
		Description: "The storage configuration of the repository docker hosted",
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"blob_store_name": {
					Description: "Blob store used to store repository contents",
					Required:    true,
					Set: func(v interface{}) int {
						return schema.HashString(strings.ToLower(v.(string)))
					},
					Type: schema.TypeString,
				},
				"strict_content_type_validation": {
					Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"write_policy": {
					Description: "Controls if deployments of and updates to assets are allowed",
					Default:     "ALLOW",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"ALLOW",
						"ALLOW_ONCE",
						"DENY",
					}, false),
				},
				"latest_policy": {
					Description: "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
					Optional:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}
	DataSourceDockerHostedStorage = &schema.Schema{
		Description: "The storage configuration of the repository docker hosted",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"blob_store_name": {
					Description: "Blob store used to store repository contents",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"strict_content_type_validation": {
					Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"write_policy": {
					Description: "Controls if deployments of and updates to assets are allowed",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"latest_policy": {
					Description: "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
					Computed:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}
)
