package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryTerraformHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing hosted terraform repository.",

		Read: dataSourceRepositoryTerraformHostedRead,
		Schema: map[string]*schema.Schema{
			"id":        common.DataSourceID,
			"name":      repositorySchema.DataSourceName,
			"online":    repositorySchema.DataSourceOnline,
			"cleanup":   repositorySchema.DataSourceCleanup,
			"component": repositorySchema.DataSourceComponent,
			"storage":   repositorySchema.DataSourceHostedStorage,
		},
	}
}

func dataSourceRepositoryTerraformHostedRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceTerraformHostedRepositoryRead(d, m)
}
