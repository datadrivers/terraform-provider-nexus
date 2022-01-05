package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing yum repository.",

		Read: dataSourceRepositoryYumHostedRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Description: "Whether this repository accepts incoming requests",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"repodata_depth": {
				Description: "Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"deploy_policy": {
				Description: "Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cleanup": getDataSourceCleanupSchema(),
			"storage": getDataSourceHostedStorageSchema(),
		},
	}
}

func dataSourceRepositoryYumHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceYumHostedRepositoryRead(d, m)
}
