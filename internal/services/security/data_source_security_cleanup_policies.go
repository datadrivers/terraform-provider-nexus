package security

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryCleanupPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a cleanup policy in Nexus Repository.",

		Read: dataSourceRepositoryCleanupPoliciesRead,

		Schema: map[string]*schema.Schema{
			"notes": {
				Description: "any details on the specific cleanup policy",
				Type:        schema.TypeString,
				Default:     "string",
				Optional:    true,
			},
			"criterialLastBlobUpdated": {
				Description: "the age of the component in days",
				Type:        schema.TypeInt,
				Default:     0,
				Optional:    false,
			},
			"criteriaLastDownloaded": {
				Description: "the last time the component had been downloaded in days",
				Type:        schema.TypeInt,
				Default:     0,
				Optional:    false,
			},
			"criteriaReleaseType": {
				Description: "When needed, this is either PRELEASE or RELEASE",
				Type:        schema.TypeString,
				Default:     "RELEASES",
				Optional:    false,
			},
			"criteriaAssetRegex": {
				Description: "a regex string to filter for specific asset paths",
				Type:        schema.TypeString,
				Default:     "string",
				Optional:    false,
			},
			"retain": {
				Description: "number of versions to keep. Only available for Docker and Maven release repositories on PostgreSQL deployments",
				Type:        schema.TypeInt,
				Default:     nil,
				Optional:    false,
			},
			"name": {
				Description: "the name of the policy needs to be unique and cannot be edited once set. Only letters, digits, underscores(_), hyphens(-), and dots(.) are allowed and may not start with underscore or dot",
				Type:        schema.TypeString,
				Default:     "policy-name",
				Optional:    true,
			},
			"format": {
				Description: "the target format for the cleanup policy. Some formats have various capabilities and requirements. Note that you cannot currently specify all formats",
				Type:        schema.TypeString,
				Default:     "string",
				Optional:    true,
			},
		},
	}
}

func dataSourceRepositoryCleanupPoliciesRead(d *schema.ResourceData, m interface{}) error {
	return resourceSecurityCleanupPoliciesRead(d, m)
}
