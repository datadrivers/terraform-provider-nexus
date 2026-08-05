package cleanuppolicies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryCleanupPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a cleanup policy in Nexus Repository.",

		Read: dataSourceRepositoryCleanupPoliciesRead,

		Schema: map[string]*schema.Schema{
			"notes": {
				Description: "Any details on the specific cleanup policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"criteria_last_blob_updated": {
				Description: "The age of the component in days",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"criteria_last_downloaded": {
				Description: "the last time the component had been downloaded in days",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"criteria_release_type": {
				Description: "When needed, this is either PRELEASE or RELEASE",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"criteria_asset_regex": {
				Description: "A regex string to filter for specific asset paths",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"retain": {
				Description: "Number of versions to keep. Only available for Docker and Maven release repositories on PostgreSQL deployments",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "The name of the policy needs to be unique and cannot be edited once set. Only letters, digits, underscores(_), hyphens(-), and dots(.) are allowed and may not start with underscore or dot",
				Type:        schema.TypeString,
				Required:    true,
			},
			"format": {
				Description: "The target format for the cleanup policy. Some formats have various capabilities and requirements. Note that you cannot currently specify all formats",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceRepositoryCleanupPoliciesRead(d *schema.ResourceData, m interface{}) error {
	return resourceCleanupPoliciesRead(d, m)
}
