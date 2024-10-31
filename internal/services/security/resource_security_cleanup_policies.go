package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityCleanupPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a cleanup policy in Nexus Repository.",

		Create: resourceSecurityCleanupPoliciesCreate,
		Read:   resourceSecurityCleanupPoliciesRead,
		Update: resourceSecurityCleanupPoliciesUpdate,
		Delete: resourceSecurityCleanupPoliciesDelete,
		Exists: resourceSecurityCleanupPoliciesExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Default:     "RELEASE",
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

func getCleanupPolicyFromResourceData(d *schema.ResourceData) security.CleanupPolicy {

	notes, _ := d.Get("notes").(string)
	criteriaLastBlobUpdated, _ := d.Get("criteria_last_blob_updated").(int)
	criteriaLastDownloaded, _ := d.Get("criteria_last_downloaded").(int)
	criteriaReleaseType, _ := d.Get("criteria_release_type").(string)
	criteriaAssetRegex, _ := d.Get("criteria_asset_regex").(string)
	retain, _ := d.Get("retain").(int)
	name, _ := d.Get("name").(string)
	format, _ := d.Get("format").(string)

	policy := security.CleanupPolicy{
		Notes:                   &notes,
		CriteriaLastBlobUpdated: &criteriaLastBlobUpdated,
		CriteriaLastDownloaded:  &criteriaLastDownloaded,
		CriteriaReleaseType:     &criteriaReleaseType,
		CriteriaAssetRegex:      &criteriaAssetRegex,
		Retain:                  retain,
		Name:                    name,
		Format:                  format,
	}

	return policy
}

func setCleanupPolicyToResourceData(cleanupPolicy *security.CleanupPolicy, d *schema.ResourceData) error {
	d.SetId(cleanupPolicy.Name)
	d.Set("notes", cleanupPolicy.Notes)
	d.Set("criteria_last_blob_updated", cleanupPolicy.CriteriaLastBlobUpdated)
	d.Set("criteria_last_downloaded", cleanupPolicy.CriteriaLastDownloaded)
	d.Set("criteria_release_type", cleanupPolicy.CriteriaReleaseType)
	d.Set("criteria_asset_regex", cleanupPolicy.CriteriaAssetRegex)
	d.Set("retain", cleanupPolicy.Retain)
	d.Set("name", cleanupPolicy.Name)
	d.Set("format", cleanupPolicy.Format)
	return nil
}

func resourceSecurityCleanupPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy := getCleanupPolicyFromResourceData(d)

	if err := client.Security.CleanupPolicy.Create(&cleanupPolicy); err != nil {
		return err
	}

	d.SetId(cleanupPolicy.Name)

	return resourceSecurityContentSelectorRead(d, m)
}

func resourceSecurityCleanupPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy, err := client.Security.CleanupPolicy.Get(d.Id())
	if err != nil {
		return err
	}

	if cleanupPolicy == nil {
		d.SetId("")
		return nil
	}

	return setCleanupPolicyToResourceData(cleanupPolicy, d)
}

func resourceSecurityCleanupPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy := getCleanupPolicyFromResourceData(d)
	if err := client.Security.CleanupPolicy.Update(&cleanupPolicy); err != nil {
		return err
	}

	return resourceSecurityCleanupPoliciesRead(d, m)
}

func resourceSecurityCleanupPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.CleanupPolicy.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityCleanupPoliciesExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	cleanupPolicy, err := client.Security.CleanupPolicy.Get(d.Id())
	return cleanupPolicy != nil, err
}
