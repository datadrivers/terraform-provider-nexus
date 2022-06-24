package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryHelmHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted yum repository.",

		Read: dataSourceRepositoryHelmHostedRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Hosted schemas
			"cleanup":   repository.DataSourceCleanup,
			"component": repository.DataSourceComponent,
			"storage":   repository.DataSourceHostedStorage,
		},
	}
}

func dataSourceRepositoryHelmHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceHelmHostedRepositoryRead(d, m)
}
