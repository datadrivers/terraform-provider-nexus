package repository

import (
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getResourceNegativeCacheSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Configuration of the negative cache handling",
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Default:     false,
					Description: "Whether to cache responses for content not present in the proxied repository",
					Optional:    true,
					Type:        schema.TypeBool,
				},
				"ttl": {
					Default:     1440,
					Description: "How long to cache the fact that a file was not found in the repository (in minutes)",
					Optional:    true,
					Type:        schema.TypeInt,
				},
			},
		},
	}
}

func flattenRepositoryNegativeCache(negativeCache *repository.NegativeCache) []map[string]interface{} {
	if negativeCache == nil {
		return nil
	}
	data := map[string]interface{}{
		"enabled": negativeCache.Enabled,
		"ttl":     negativeCache.TTL,
	}
	return []map[string]interface{}{data}
}
