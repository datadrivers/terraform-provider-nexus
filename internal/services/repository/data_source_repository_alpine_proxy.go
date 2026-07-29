package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryAlpineProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing alpine proxy repository.",

		Read: dataSourceRepositoryAlpineProxyRead,

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
			// Alpine signing
			"alpine_signing": {
				Description: "PGP signing key for the alpine repository",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair": {
							Description: "PEM-encoded RSA private key used to sign APKINDEX",
							Computed:    true,
							Sensitive:   true,
							Type:        schema.TypeString,
						},
						"passphrase": {
							Description: "Passphrase to access the signing key",
							Computed:    true,
							Sensitive:   true,
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryAlpineProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceAlpineProxyRepositoryRead(resourceData, m)
}
