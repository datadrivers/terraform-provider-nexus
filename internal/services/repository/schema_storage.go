package repository

import (
	"strings"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getResourceHostedStorageSchema() *schema.Schema {
	return &schema.Schema{
		DefaultFunc: repositoryStorageDefault,
		Description: "The storage configuration of the repository",
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"blob_store_name": {
					Default:     "default",
					Description: "Blob store used to store repository contents",
					Optional:    true,
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
}

func getDataSourceHostedStorageSchema() *schema.Schema {
	return &schema.Schema{
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
}

func repositoryStorageDefault() (interface{}, error) {

	data := map[string]interface{}{
		"blob_store_name":                "default",
		"strict_content_type_validation": true,
	}
	return []map[string]interface{}{data}, nil
}

func flattenStorage(storage *repository.Storage, d *schema.ResourceData) []map[string]interface{} {
	if storage == nil {
		return nil
	}
	data := map[string]interface{}{
		"blob_store_name":                storage.BlobStoreName,
		"strict_content_type_validation": storage.StrictContentTypeValidation,
	}
	return []map[string]interface{}{data}
}

func flattenHostedStorage(storage *repository.HostedStorage, d *schema.ResourceData) []map[string]interface{} {
	if storage == nil {
		return nil
	}
	data := map[string]interface{}{
		"blob_store_name":                storage.BlobStoreName,
		"strict_content_type_validation": storage.StrictContentTypeValidation,
	}
	if storage.WritePolicy != nil {
		data["write_policy"] = storage.WritePolicy
	}
	return []map[string]interface{}{data}
}
