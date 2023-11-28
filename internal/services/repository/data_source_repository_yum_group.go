package repository

import (
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryYumGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing yum group repository.",

		Read: dataSourceRepositoryYumGroupRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Group schemas
			"group":   repository.DataSourceGroup,
			"storage": repository.DataSourceStorage,
			// Yum hosted schemas
			"yum_signing": repository.DataSourceYumSigning,
		},
	}
}

func dataSourceRepositoryYumGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceYumGroupRepositoryRead(resourceData, m)
}
