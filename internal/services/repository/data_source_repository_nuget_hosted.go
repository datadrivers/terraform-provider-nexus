package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/common"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryNugetHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted nuget repository.",

		Read: dataSourceRepositoryNugetHostedRead,
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

func dataSourceRepositoryNugetHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceNugetHostedRepositoryRead(d, m)
}
