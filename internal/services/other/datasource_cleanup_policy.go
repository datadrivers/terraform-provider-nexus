package other

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceCleanupPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to work with cleanup policy.",

		Read: dataSourceCleanupPolicyRead,
		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"name": {
				Description:  "Use a unique name for the cleanup policy",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"format": {
				Description: "The format that this cleanup policy can be applied to",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"retain": {
				Description: "number of versions to keep. Require: Nexus Pro",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"criteria": {
				Description: "Remove all components that match all selected criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_blob_updated": {
							Description: "The age of the component in days.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"last_downloaded": {
							Description: "The last time the component had been downloaded in days.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"release_type": {
							Description: "Is one of: RELEASES_AND_PRERELEASES, PRERELEASES, RELEASES], Only maven2, npm and yum repositories support this field.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"asset_regex": {
							Description: "A regex string to filter for specific asset paths. Not for gitlfs or *",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCleanupPolicyRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceCleanupPolicyRead(d, m)
}
