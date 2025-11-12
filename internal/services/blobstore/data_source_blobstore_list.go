package blobstore

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlobstoreList() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a list with all Blob Stores.",

		Read: dataSourceBlobstoreList,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"items": {
				Description: "A List of all Blob Stores",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": blobstore.DataSourceName,
						"type": {
							Description: "The type of current blob store",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBlobstoreList(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	items := []map[string]string{}
	bs, err := nexusClient.BlobStore.List()
	if err != nil {
		return err
	}

	for _, blobstore := range bs {
		items = append(items, map[string]string{
			"name": blobstore.Name,
			"type": blobstore.Type,
		})
	}

	if err := resourceData.Set("items", items); err != nil {
		return err
	}
	resourceData.SetId("blobStoreList")
	return nil
}
