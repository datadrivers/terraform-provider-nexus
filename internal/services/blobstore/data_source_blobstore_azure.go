package blobstore

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlobstoreAzure() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this data source to get details of an existing Nexus Azure blobstore.`,

		Read: dataSourceBlobstoreAzureRead,
		Schema: map[string]*schema.Schema{
			"id":                  common.DataSourceID,
			"name":                blobstore.DataSourceName,
			"blob_count":          blobstore.DataSourceBlobCount,
			"soft_quota":          blobstore.DataSourceSoftQuota,
			"total_size_in_bytes": blobstore.DataSourceTotalSizeInBytes,
			"bucket_configuration": {
				Description: "The Azure specific configuration details for the Azure object that'll contain the blob store",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Description: "Account name found under Access keys for the storage account",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"authentication": {
							Description: "The Azure specific authentication details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_method": {
										Description: "The type of Azure authentication to use. Possible values: `ACCOUNTKEY` and `MANAGEDIDENTITY`",
										Computed:    true,
										Type:        schema.TypeString,
									},
								},
							},
							Computed: true,
							Type:     schema.TypeList,
						},
						"container_name": {
							Description: "The name of an existing container to be used for storage",
							Computed:    true,
							Type:        schema.TypeString,
						},
					},
				},
				Computed: true,
				Type:     schema.TypeList,
			},
		},
	}
}

func dataSourceBlobstoreAzureRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBlobstoreAzureRead(resourceData, m)
}
