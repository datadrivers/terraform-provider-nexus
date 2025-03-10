package cleanup_policy

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCleanupPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCleanupPolicyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"criteria_last_blob_updated": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"criteria_last_downloaded": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"criteria_release_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"criteria_asset_regex": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"format": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCleanupPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	id := d.Get("id").(string)

	policy, err := client.Cleanup.Get(id)
	if err != nil {
		return err
	}

	//d.SetId(policy.ID)
	d.Set("name", policy.Name)
	d.Set("notes", policy.Notes)
	d.Set("criteria_last_blob_updated", policy.CriteriaLastBlobUpdated)
	d.Set("criteria_last_downloaded", policy.CriteriaLastDownloaded)
	d.Set("criteria_release_type", policy.CriteriaReleaseType)
	d.Set("criteria_asset_regex", policy.CriteriaAssetRegex)
	d.Set("retain", policy.Retain)
	d.Set("format", policy.Format)

	return nil
}
