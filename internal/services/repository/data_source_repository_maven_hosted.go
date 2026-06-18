package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/common"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryMavenHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted yum repository.",

		Read: dataSourceRepositoryMavenHostedRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Hosted schemas
			"cleanup":   repository.DataSourceCleanup,
			"component": repository.DataSourceComponent,
			"storage":   repository.DataSourceHostedStorage,
			// Maven hosted schemas
			"maven": repository.DataSourceMaven,
		},
	}
}

func dataSourceRepositoryMavenHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceMavenHostedRepositoryRead(d, m)
}
