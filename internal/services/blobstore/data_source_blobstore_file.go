package blobstore

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlobstoreFile() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get details of an existing Nexus File blobstore.",

		Read: dataSourceBlobstoreFileRead,
		Schema: map[string]*schema.Schema{
			"id":   common.DataSourceID,
			"name": blobstore.DataSourceName,
			"path": {
				Description: "The path to the blobstore contents",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"available_space_in_bytes": blobstore.DataSourceAvailableSpaceInBytes,
			"blob_count":               blobstore.DataSourceBlobCount,
			"soft_quota":               blobstore.DataSourceSoftQuota,
			"total_size_in_bytes":      blobstore.DataSourceTotalSizeInBytes,
		},
	}
}

func dataSourceBlobstoreFileRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBlobstoreFileRead(resourceData, m)
}
