package repository

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	ResourceStorage = &schema.Schema{
		Description: "The storage configuration of the repository",
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
					Default:     true,
					Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
					Optional:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}
	DataSourceStorage = &schema.Schema{
		Description: "The storage configuration of the repository",
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
			},
		},
	}

	ResourceHostedStorage = &schema.Schema{
		Description: "The storage configuration of the repository",
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
			},
		},
	}
	DataSourceHostedStorage = &schema.Schema{
		Description: "The storage configuration of the repository",
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
			},
		},
	}
)
