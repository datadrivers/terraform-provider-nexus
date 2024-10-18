package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryNpmProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing npm proxy repository.",

		Read: dataSourceRepositoryNpmProxyRead,

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repositorySchema.DataSourceName,
			"online": repositorySchema.DataSourceOnline,
			// Proxy schemas
			"cleanup":        repositorySchema.DataSourceCleanup,
			"http_client":    repositorySchema.DataSourceHTTPClient,
			"negative_cache": repositorySchema.DataSourceNegativeCache,
			"proxy":          repositorySchema.DataSourceProxy,
			"routing_rule":   repositorySchema.DataSourceRoutingRule,
			"storage":        repositorySchema.DataSourceStorage,
			// NPM proxy schemas
			"remove_non_cataloged": {
				Description: "Remove non-catalogued versions from the npm package metadata.",
				Computed:    true,
				Type:        schema.TypeBool,
				Deprecated:  "This field is removed in nexus version 3.66.0 and will be removed in the next major release of this provider",
			},
			"remove_quarantined": {
				Description: "Remove quarantined versions from the npm package metadata.",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func dataSourceRepositoryNpmProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceNpmProxyRepositoryRead(resourceData, m)
}
