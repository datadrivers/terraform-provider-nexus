package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/common"
	"github.com/williamt1996/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryRubygemsGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing rubygems group repository.",

		Read: dataSourceRepositoryRubygemsGroupRead,
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

func dataSourceRepositoryRubygemsGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceRubygemsGroupRepositoryRead(resourceData, m)
}
