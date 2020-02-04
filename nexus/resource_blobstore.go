package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBlobstore() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlobstoreCreate,
		Read:   resourceBlobstoreRead,
		Update: resourceBlobstoreUpdate,
		Delete: resourceBlobstoreDelete,

		Schema: map[string]*schema.Schema{
			"available_space_in_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"blob_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"soft_quota": {
				Description: "The limit in MB.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Default:  0,
							Optional: true,
							Type:     schema.TypeInt,
						},
						"type": {
							Description:  "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Optional:     true,
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
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getBlobstoreFromResourceData(d *schema.ResourceData) nexus.Blobstore {
	bs := nexus.Blobstore{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		Path: d.Get("path").(string),
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
	// Path not returned by API :-/
	//d.Set("path", bs.Path)
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
	nexusClient := m.(nexus.Client)
	id := d.Id()

	if d.HasChange("path") {
		bs := getBlobstoreFromResourceData(d)
		if err := nexusClient.BlobstoreUpdate(id, bs); err != nil {
			return err
		}
	}
	return nil
}

func resourceBlobstoreDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	id := d.Id()

	if err := nexusClient.BlobstoreDelete(id); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
