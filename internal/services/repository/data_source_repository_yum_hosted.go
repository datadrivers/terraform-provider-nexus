package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted yum repository.",

		Read: dataSourceRepositoryYumHostedRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Hosted schemas
			"cleanup":   repository.DataSourceCleanup,
			"component": repository.DataSourceComponent,
			"storage":   repository.DataSourceHostedStorage,
			// Yum hosted schemas
			"deploy_policy": {
				Description: "Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"repodata_depth": {
				Description: "Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceRepositoryYumHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceYumHostedRepositoryRead(d, m)
}
