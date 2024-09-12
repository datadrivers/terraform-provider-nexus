package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryDockerProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing docker proxy repository.",

		Read: dataSourceRepositoryDockerProxyRead,

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
			// Docker proxy schemas
			"docker": repositorySchema.DataSourceDocker,
			"docker_proxy": {
				Description: "docker_proxy contains the configuration of the docker index",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_type": {
							Description: "Type of Docker Index",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"index_url": {
							Description: "Url of Docker Index to use",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"cache_foreign_layers": {
							Description: "Allow Nexus Repository Manager to download and cache foreign layers",
							Computed:    true,
							Type:        schema.TypeBool,
						},
						"foreign_layer_url_whitelist": {
							Description: "A set of regular expressions used to identify URLs that are allowed for foreign layer requests",
							Computed:    true,
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryDockerProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceDockerProxyRepositoryRead(resourceData, m)
}
