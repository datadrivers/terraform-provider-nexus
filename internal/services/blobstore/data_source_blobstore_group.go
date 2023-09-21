package blobstore

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlobstoreGroup() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this data source to get details of an existing Nexus Group blobstore.`,

		Read: dataSourceBlobstoreGroupRead,
		Schema: map[string]*schema.Schema{
			"id":                       common.DataSourceID,
			"name":                     blobstore.DataSourceName,
			"available_space_in_bytes": blobstore.ResourceAvailableSpaceInBytes,
			"blob_count":               blobstore.DataSourceBlobCount,
			"fill_policy": {
				Description: "The policy how to fill the members. Possible values: `roundRobin` or `writeToFirst`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"members": {
				Description: "List of the names of blob stores that are members of this group",
				Computed:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"soft_quota":          blobstore.DataSourceSoftQuota,
			"total_size_in_bytes": blobstore.DataSourceTotalSizeInBytes,
		},
	}
}

func dataSourceBlobstoreGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBlobstoreGroupRead(resourceData, m)
}
