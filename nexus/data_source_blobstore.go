package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBlobstore() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlobstoreRead,

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
			"soft_quota": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"total_size_in_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceBlobstoreRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceBlobstoreRead(d, m)
}
