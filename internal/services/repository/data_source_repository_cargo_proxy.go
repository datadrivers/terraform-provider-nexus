package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryCargoProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing cargo proxy repository.",

		Read: dataSourceRepositoryCargoProxyRead,

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
			// Cargo proxy schemas
			"cargo_version": {
				Description: "Cargo protocol version",
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

func dataSourceRepositoryCargoProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceCargoProxyRepositoryRead(resourceData, m)
}
