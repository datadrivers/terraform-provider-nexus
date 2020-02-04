package nexus

import (
	"fmt"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBlobstore() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlobstoreCreate,
		Read:   resourceBlobstoreRead,
		Update: resourceBlobstoreUpdate,
		Delete: resourceBlobstoreDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blob_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_size_in_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_space_in_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getBlobstoreFromResourceData(d *schema.ResourceData) nexus.Blobstore {
	bs := nexus.Blobstore{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	if _, ok := d.GetOk("soft_quota"); ok {
		softQuotaList := d.Get("soft_quota").([]interface{})
		softQuotaConfig := softQuotaList[0].(map[string]interface{})

		bs.BlobstoreSoftQuota = &nexus.BlobstoreSoftQuota{
			Limit: softQuotaConfig["limit"].(int),
			Type:  softQuotaConfig["type"].(string),
		}
	}

	return bs
}

func resourceBlobstoreCreate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs := getBlobstoreFromResourceData(d)

	if err := nexusClient.BlobstoreCreate(bs); err != nil {
		return err
	}

	d.SetId(bs.Name)

	return nil
}

func resourceBlobstoreRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	id := d.Get("name").(string)

	bs, err := nexusClient.BlobstoreRead(id)
	if err != nil {
		return err
	}

	if bs == nil {
		d.SetId("")
		return nil
	}

	d.SetId(bs.Name)
	d.Set("available_space_in_bytes", bs.AvailableSpaceInBytes)
	d.Set("blob_count", bs.BlobCount)
	d.Set("name", bs.Name)
	d.Set("path", bs.Path)
	if bs.BlobstoreSoftQuota != nil {
		if err := d.Set("soft_quota", flattenBlobstoreSoftQuota(bs.BlobstoreSoftQuota)); err != nil {
			return err
		}
	}
	d.Set("total_size_in_bytes", bs.TotalSizeInBytes)
	d.Set("type", bs.Type)

	return nil
}

func flattenBlobstoreSoftQuota(softQuota *nexus.BlobstoreSoftQuota) []map[string]interface{} {
	data := map[string]interface{}{
		"limit": softQuota.Limit,
		"type":  softQuota.Type,
	}
	return []map[string]interface{}{data}
}

func resourceBlobstoreUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("resourceBlobstoreUpdate not yet implemented")
}

func resourceBlobstoreDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("resourceBlobstoreDelete not yet implemented")
}
