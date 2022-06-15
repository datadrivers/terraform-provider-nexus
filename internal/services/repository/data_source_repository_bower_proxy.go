package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryBowerProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing bower proxy repository.",

		Read: dataSourceRepositoryBowerProxyRead,

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
			// Bower proxy schemas
			"rewrite_package_urls": {
				Description: "Whether to force Bower to retrieve packages through this proxy repository",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func dataSourceRepositoryBowerProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBowerProxyRepositoryRead(resourceData, m)
}
