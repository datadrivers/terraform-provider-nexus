package blobstore

import (
	"fmt"
	"log"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	blobstoreSchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBlobstoreGroup() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this resource to create a Nexus group blobstore.`,

		Create: resourceBlobstoreGroupCreate,
		Read:   resourceBlobstoreGroupRead,
		Update: resourceBlobstoreGroupUpdate,
		Delete: resourceBlobstoreGroupDelete,
		Exists: resourceBlobstoreGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id":                       common.ResourceID,
			"name":                     blobstoreSchema.ResourceName,
			"available_space_in_bytes": blobstoreSchema.ResourceAvailableSpaceInBytes,
			"blob_count":               blobstoreSchema.ResourceBlobCount,
			"fill_policy": {
				Description:  "The policy how to fill the members. Possible values: `roundRobin` or `writeToFirst`",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(blobstore.GroupFillPolicyRoundRobin), string(blobstore.GroupFillPolicyWriteToFirst)}, false),
				Required:     true,
			},
			"members": {
				Description: "List of the names of blob stores that are members of this group",
				Required:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
			},
			"soft_quota":          blobstoreSchema.ResourceSoftQuota,
			"total_size_in_bytes": blobstoreSchema.ResourceTotalSizeInBytes,
		},
	}
}

func getBlobstoreGroupFromResourceData(resourceData *schema.ResourceData) blobstore.Group {
	bs := blobstore.Group{
		Name:       resourceData.Get("name").(string),
		FillPolicy: resourceData.Get("fill_policy").(string),
		Members:    tools.ConvertStringSet(resourceData.Get("members").(*schema.Set)),
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

func resourceBlobstoreGroupCreate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreGroupFromResourceData(resourceData)

	if err := nexusClient.BlobStore.Group.Create(&bs); err != nil {
		return err
	}

	resourceData.SetId(bs.Name)
	err := resourceData.Set("name", bs.Name)
	if err != nil {
		return err
	}

	return resourceBlobstoreGroupRead(resourceData, m)
}

func resourceBlobstoreGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.Group.Get(resourceData.Id())
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
	if err := resourceData.Set("fill_policy", string(bs.FillPolicy)); err != nil {
		return err
	}
	if err := resourceData.Set("members", bs.Members); err != nil {
		return err
	}
	if err := resourceData.Set("name", bs.Name); err != nil {
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

func resourceBlobstoreGroupUpdate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreGroupFromResourceData(resourceData)
	if err := nexusClient.BlobStore.Group.Update(resourceData.Id(), &bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreGroupDelete(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	if err := nexusClient.BlobStore.Group.Delete(resourceData.Id()); err != nil {
		return err
	}

	resourceData.SetId("")

	return nil
}

func resourceBlobstoreGroupExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.Group.Get(resourceData.Id())
	return bs != nil, err
}
