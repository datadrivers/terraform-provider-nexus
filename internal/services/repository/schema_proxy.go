package repository

import (
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getResourceProxySchema() *schema.Schema {
	return &schema.Schema{
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
}

func flattenRepositoryProxy(proxy *repository.Proxy) []map[string]interface{} {
	if proxy == nil {
		return nil
	}
	data := map[string]interface{}{
		"content_max_age":  proxy.ContentMaxAge,
		"metadata_max_age": proxy.MetadataMaxAge,
		"remote_url":       proxy.RemoteURL,
	}
	return []map[string]interface{}{data}
}
