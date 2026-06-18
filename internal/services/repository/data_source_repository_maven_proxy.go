package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/williamt1996/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryMavenProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing maven proxy repository.",

		Read: dataSourceRepositoryMavenProxyRead,

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repositorySchema.DataSourceName,
			"online": repositorySchema.DataSourceOnline,
			// Proxy schemas
			"cleanup":        repositorySchema.DataSourceCleanup,
			"http_client":    repositorySchema.DataSourceHTTPClientWithPreemptiveAuth,
			"negative_cache": repositorySchema.DataSourceNegativeCache,
			"proxy":          repositorySchema.DataSourceProxy,
			"routing_rule":   repositorySchema.DataSourceRoutingRule,
			"storage":        repositorySchema.DataSourceStorage,
			// Maven proxy schemas
			"maven": repositorySchema.DataSourceMaven,
		},
	}
}

func dataSourceRepositoryMavenProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceMavenProxyRepositoryRead(resourceData, m)
}
