package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/williamt1997/terraform-provider-nexus/internal/schema/common"
	"github.com/williamt1997/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryPypiGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing pypi group repository.",

		Read: dataSourceRepositoryPypiGroupRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Group schemas
			"group":   repository.DataSourceGroupDeploy,
			"storage": repository.DataSourceStorage,
		},
	}
}

func dataSourceRepositoryPypiGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourcePypiGroupRepositoryRead(resourceData, m)
}
