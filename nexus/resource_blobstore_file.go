/*
Use this resource to create a Nexus blobstore.

Example Usage

Example Usage for file store

```hcl
resource "nexus_blobstore_file" "default" {
  name = "blobstore-file"
  type = "File"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}
```

*/
package nexus

import (
	"fmt"
	"log"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBlobstoreFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlobstoreFileCreate,
		Read:   resourceBlobstoreFileRead,
		Update: resourceBlobstoreFileUpdate,
		Delete: resourceBlobstoreFileDelete,
		Exists: resourceBlobstoreFileExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Blobstore name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				ConflictsWith: []string{"bucket_configuration"},
				Description:   "The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory",
				Type:          schema.TypeString,
				Optional:      true,
			},
			"available_space_in_bytes": {
				Description: "Available space in Bytes",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"blob_count": {
				Description: "Count of blobs",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"soft_quota": {
				Description: "Soft quota of the blobstore",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Description:  "The limit in Bytes. Minimum value is 1000000",
							Required:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntAtLeast(100000),
						},
						"type": {
							Description:  "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"spaceRemainingQuota", "spaceUsedQuota"}, false),
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"total_size_in_bytes": {
				Description: "The total size of the blobstore in Bytes",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func getBlobstoreFileFromResourceData(resourceData *schema.ResourceData) nexus.Blobstore {
	bs := nexus.Blobstore{
		Name: resourceData.Get("name").(string),
		Type: "File",
	}

	bs.Path = resourceData.Get("path").(string)

	if _, ok := resourceData.GetOk("soft_quota"); ok {
		softQuotaList := resourceData.Get("soft_quota").([]interface{})
		softQuotaConfig := softQuotaList[0].(map[string]interface{})

		bs.BlobstoreSoftQuota = &nexus.BlobstoreSoftQuota{
			Limit: softQuotaConfig["limit"].(int),
			Type:  softQuotaConfig["type"].(string),
		}
	}

	return bs
}

func resourceBlobstoreFileCreate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs := getBlobstoreFileFromResourceData(resourceData)

	if err := nexusClient.BlobstoreCreate(bs); err != nil {
		return err
	}

	resourceData.SetId(bs.Name)
	resourceData.Set("name", bs.Name)

	return resourceBlobstoreRead(resourceData, m)
}

func resourceBlobstoreFileRead(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs, err := nexusClient.BlobstoreRead(resourceData.Id())
	log.Print(bs)
	if err != nil {
		return err
	}

	if bs == nil {
		resourceData.SetId("")
		return nil
	}

	resourceData.Set("available_space_in_bytes", bs.AvailableSpaceInBytes)
	resourceData.Set("blob_count", bs.BlobCount)
	resourceData.Set("name", bs.Name)
	resourceData.Set("path", bs.Path)
	resourceData.Set("total_size_in_bytes", bs.TotalSizeInBytes)

	if bs.BlobstoreSoftQuota != nil {
		if err := resourceData.Set("soft_quota", flattenBlobstoreFileSoftQuota(bs.BlobstoreSoftQuota)); err != nil {
			return fmt.Errorf("Error reading soft quota: %s", err)
		}
	}

	return nil
}

func resourceBlobstoreFileUpdate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs := getBlobstoreFileFromResourceData(resourceData)
	if err := nexusClient.BlobstoreUpdate(resourceData.Id(), bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreFileDelete(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if err := nexusClient.BlobstoreDelete(resourceData.Id()); err != nil {
		return err
	}

	resourceData.SetId("")

	return nil
}

func resourceBlobstoreFileExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	bs, err := nexusClient.BlobstoreRead(resourceData.Id())
	return bs != nil, err
}

func flattenBlobstoreFileSoftQuota(softQuota *nexus.BlobstoreSoftQuota) []map[string]interface{} {
	if softQuota == nil {
		return nil
	}
	data := map[string]interface{}{
		"limit": softQuota.Limit,
		"type":  softQuota.Type,
	}
	return []map[string]interface{}{data}
}
