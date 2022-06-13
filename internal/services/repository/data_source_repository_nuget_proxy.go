package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryNugetProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing nuget proxy repository.",

		Read: dataSourceRepositoryNugetProxyRead,

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
			// Nuget proxy schemas
			"nuget_version": {
				Description: "Nuget protocol version",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"query_cache_item_max_age": {
				Description: "How long to cache query results from the proxied repository (in seconds)",
				Computed:    true,
				Type:        schema.TypeInt,
			},
		},
	}
}

func dataSourceRepositoryNugetProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceNugetProxyRepositoryRead(resourceData, m)
}
