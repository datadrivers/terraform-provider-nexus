package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryPypiHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted pypi repository.",

		Read: dataSourceRepositoryPypiHostedRead,
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

func dataSourceRepositoryPypiHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourcePypiHostedRepositoryRead(d, m)
}
