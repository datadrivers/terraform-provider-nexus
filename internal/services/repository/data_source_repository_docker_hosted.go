package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryDockerHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing docker repository.",

		Read: dataSourceRepositoryDockerHostedRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Hosted schemas
			"cleanup":   repository.DataSourceCleanup,
			"component": repository.DataSourceComponent,
			"storage":   repository.DataSourceHostedStorage,
			// Docker hosted schemas
			"docker": repository.DataSourceDocker,
		},
	}
}

func dataSourceRepositoryDockerHostedRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceDockerHostedRepositoryRead(resourceData, m)
}
