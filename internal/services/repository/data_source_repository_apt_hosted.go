package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryAptHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing apt repository.",

		Read: dataSourceRepositoryAptHostedRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Hosted schemas
			"cleanup":   repository.DataSourceCleanup,
			"component": repository.DataSourceComponent,
			"storage":   repository.DataSourceHostedStorage,
			// Apt hosted schemas
			"distribution": {
				Description: "Distribution to fetch",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"signing": {
				Description: "Contains signing data of repositores",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair": {
							Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"passphrase": {
							Description: "Passphrase to access PGP signing key.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryAptHostedRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceAptHostedRepositoryRead(resourceData, m)
}
