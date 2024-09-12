package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryBowerGroup() *schema.Resource {
	return &schema.Resource{
		Description: `!> This data source is deprecated and will be removed in the next major release of this provider. Bower repositories were removed in Nexus 3.71.0.

Use this data source to get an existing bower group repository.`,
		DeprecationMessage: "This data source is deprecated and will be removed in the next major release of this provider. Bower repositories were removed in Nexus 3.71.0.",

		Read: dataSourceRepositoryBowerGroupRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Group schemas
			"group":   repository.DataSourceGroup,
			"storage": repository.DataSourceStorage,
		},
	}
}

func dataSourceRepositoryBowerGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBowerGroupRepositoryRead(resourceData, m)
}
