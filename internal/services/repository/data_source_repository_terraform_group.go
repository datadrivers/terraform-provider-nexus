package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryTerraformGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing group terraform repository.",

		Read: dataSourceRepositoryTerraformGroupRead,
		Schema: map[string]*schema.Schema{
			"id":      common.DataSourceID,
			"name":    repositorySchema.DataSourceName,
			"online":  repositorySchema.DataSourceOnline,
			"group":   repositorySchema.DataSourceGroup,
			"storage": repositorySchema.DataSourceStorage,
		},
	}
}

func dataSourceRepositoryTerraformGroupRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceTerraformGroupRepositoryRead(d, m)
}
