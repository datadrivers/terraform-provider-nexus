package repository

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/dre2004/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryAptProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing apt proxy repository.",

		Read: dataSourceRepositoryAptProxyRead,

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
			// Apt proxy schemas
			"distribution": {
				Description: "Distribution to fetch",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"flat": {
				Description: "Distribution to fetch",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func dataSourceRepositoryAptProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceAptProxyRepositoryRead(resourceData, m)
}
