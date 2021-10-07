/*
Use this to create a Nexus blobstore.

Example Usage

```hcl
data "nexus_blobstore" "docker" {
	name = "docker"
}
```
*/
package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional:    true,
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
							Optional:    true,
						},
						"type": {
							Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
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

	return resourceBlobstoreRead(resourceData, m)
}
