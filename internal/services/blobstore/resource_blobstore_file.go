package blobstore

import (
	"fmt"
	"log"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	blobstoreSchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBlobstoreFile() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus file blobstore.",

		Create: resourceBlobstoreFileCreate,
		Read:   resourceBlobstoreFileRead,
		Update: resourceBlobstoreFileUpdate,
		Delete: resourceBlobstoreFileDelete,
		Exists: resourceBlobstoreFileExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id":   common.ResourceID,
			"name": blobstoreSchema.ResourceName,
			"path": {
				Description: "The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"available_space_in_bytes": blobstoreSchema.ResourceAvailableSpaceInBytes,
			"blob_count":               blobstoreSchema.ResourceBlobCount,
			"soft_quota":               blobstoreSchema.ResourceSoftQuota,
			"total_size_in_bytes":      blobstoreSchema.ResourceTotalSizeInBytes,
		},
	}
}

func getBlobstoreFileFromResourceData(resourceData *schema.ResourceData) blobstore.File {
	bs := blobstore.File{
		Name: resourceData.Get("name").(string),
	}

	if _, ok := resourceData.GetOk("path"); ok {
		bs.Path = resourceData.Get("path").(string)
	}

	if _, ok := resourceData.GetOk("soft_quota"); ok {
		softQuotaConfig := resourceData.Get("soft_quota").([]interface{})[0].(map[string]interface{})

		bs.SoftQuota = &blobstore.SoftQuota{
			Limit: int64(softQuotaConfig["limit"].(int)),
			Type:  softQuotaConfig["type"].(string),
		}
	}

	return bs
}

func resourceBlobstoreFileCreate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreFileFromResourceData(resourceData)

	if err := nexusClient.BlobStore.File.Create(&bs); err != nil {
		return err
	}

	resourceData.SetId(bs.Name)
	err := resourceData.Set("name", bs.Name)
	if err != nil {
		return err
	}

	return resourceBlobstoreFileRead(resourceData, m)
}

func resourceBlobstoreFileRead(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.File.Get(resourceData.Id())
	log.Printf("[DEBUG] BlobStore:\n%+v\n", bs)
	if err != nil {
		return err
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := nexusClient.BlobStore.List()
	if err != nil {
		return err
	}
	for _, generic := range genericBlobstores {
		if generic.Name == bs.Name {
			genericBlobstoreInformation = generic
		}
	}

	if bs == nil {
		resourceData.SetId("")
		return nil
	}

	if err := resourceData.Set("available_space_in_bytes", genericBlobstoreInformation.AvailableSpaceInBytes); err != nil {
		return err
	}
	if err := resourceData.Set("blob_count", genericBlobstoreInformation.BlobCount); err != nil {
		return err
	}
	if err := resourceData.Set("name", bs.Name); err != nil {
		return err
	}
	if err := resourceData.Set("path", bs.Path); err != nil {
		return err
	}
	if err := resourceData.Set("total_size_in_bytes", genericBlobstoreInformation.TotalSizeInBytes); err != nil {
		return err
	}

	if bs.SoftQuota != nil {
		if err := resourceData.Set("soft_quota", flattenSoftQuota(bs.SoftQuota)); err != nil {
			return fmt.Errorf("error reading soft quota: %s", err)
		}
	}

	return nil
}

func resourceBlobstoreFileUpdate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreFileFromResourceData(resourceData)
	if err := nexusClient.BlobStore.File.Update(resourceData.Id(), &bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreFileDelete(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	if err := nexusClient.BlobStore.File.Delete(resourceData.Id()); err != nil {
		return err
	}

	resourceData.SetId("")

	return nil
}

func resourceBlobstoreFileExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.File.Get(resourceData.Id())
	return bs != nil, err
}
