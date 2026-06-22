package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/williamt1997/terraform-provider-nexus/internal/schema/common"
	"github.com/williamt1997/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryCargoHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted cargo repository.",

		Read: dataSourceRepositoryCargoHostedRead,
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

func dataSourceRepositoryCargoHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceCargoHostedRepositoryRead(d, m)
}
