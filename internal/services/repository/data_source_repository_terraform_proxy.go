package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryTerraformProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing proxy terraform repository.",

		Read: dataSourceRepositoryTerraformProxyRead,
		Schema: map[string]*schema.Schema{
			"id":             common.DataSourceID,
			"name":           repositorySchema.DataSourceName,
			"online":         repositorySchema.DataSourceOnline,
			"cleanup":        repositorySchema.DataSourceCleanup,
			"http_client":    repositorySchema.DataSourceHTTPClientWithPreemptiveAuth,
			"negative_cache": repositorySchema.DataSourceNegativeCache,
			"proxy":          repositorySchema.DataSourceProxy,
			"routing_rule":   repositorySchema.DataSourceRoutingRule,
			"storage":        repositorySchema.DataSourceStorage,
		},
	}
}

func dataSourceRepositoryTerraformProxyRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceTerraformProxyRepositoryRead(d, m)
}
