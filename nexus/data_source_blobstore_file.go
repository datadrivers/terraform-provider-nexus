/*
Use this to read an existing Nexus file blobstore.

Example Usage

```hcl
data "nexus_blobstore_file" "default" {
	name = "default"
}
```
*/
package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBlobstoreFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlobstoreFileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Blobstore name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"path": {
				Description: "The path to the blobstore contents",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"available_space_in_bytes": {
				Computed:    true,
				Description: "Available space in Bytes",
				Type:        schema.TypeInt,
			},
			"blob_count": {
				Computed:    true,
				Description: "Count of blobs",
				Type:        schema.TypeInt,
			},
			"soft_quota": {
				Description: "Soft quota of the blobstore",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Description: "The limit in Bytes. Minimum value is 1000000",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"type": {
							Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Computed: true,
				Type:     schema.TypeList,
			},
			"total_size_in_bytes": {
				Computed:    true,
				Description: "The total size of the blobstore in Bytes",
				Type:        schema.TypeInt,
			},
		},
	}
}

func dataSourceBlobstoreFileRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBlobstoreFileRead(resourceData, m)
}
