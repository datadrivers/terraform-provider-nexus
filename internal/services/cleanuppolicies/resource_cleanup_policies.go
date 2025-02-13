package cleanuppolicies

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCleanupPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a cleanup policy in Nexus Repository.",

		Create: resourceCleanupPoliciesCreate,
		Read:   resourceCleanupPoliciesRead,
		Update: resourceCleanupPoliciesUpdate,
		Delete: resourceCleanupPoliciesDelete,
		Exists: resourceCleanupPoliciesExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"notes": {
				Description: "Any details on the specific cleanup policy",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"criteria_last_blob_updated": {
				Description: "The age of the component in days",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"criteria_last_downloaded": {
				Description: "the last time the component had been downloaded in days",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"criteria_release_type": {
				Description: "When needed, this is either PRELEASE or RELEASE",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"criteria_asset_regex": {
				Description: "A regex string to filter for specific asset paths",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"retain": {
				Description: "Number of versions to keep. Only available for Docker and Maven release repositories on PostgreSQL deployments",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"name": {
				Description: "The name of the policy needs to be unique and cannot be edited once set. Only letters, digits, underscores(_), hyphens(-), and dots(.) are allowed and may not start with underscore or dot",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"format": {
				Description: "The target format for the cleanup policy. Some formats have various capabilities and requirements. Note that you cannot currently specify all formats",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func getCleanupPolicyFromResourceData(d *schema.ResourceData) cleanuppolicies.CleanupPolicy {

	notes, _ := d.Get("notes").(string)
	criteriaLastBlobUpdated, _ := d.Get("criteria_last_blob_updated").(int)
	criteriaLastDownloaded, _ := d.Get("criteria_last_downloaded").(int)
	criteriaReleaseType, _ := d.Get("criteria_release_type").(string)
	criteriaAssetRegex, _ := d.Get("criteria_asset_regex").(string)
	retain, _ := d.Get("retain").(int)
	name, _ := d.Get("name").(string)
	format, _ := d.Get("format").(string)

	policy := cleanuppolicies.CleanupPolicy{
		Notes:                   &notes,
		CriteriaLastBlobUpdated: &criteriaLastBlobUpdated,
		CriteriaLastDownloaded:  &criteriaLastDownloaded,
		CriteriaAssetRegex:      &criteriaAssetRegex,
		Name:                    name,
		Format:                  format,
	}

	// Only set CriteriaReleaseType for applicable formats
	if format == "maven2" || format == "npm" || format == "yum" {
		policy.CriteriaReleaseType = &criteriaReleaseType
	}

	// Only set CriteriaReleaseType for applicable formats
	if format == "maven2" || format == "docker" {
		policy.Retain = retain
	}

	return policy
}

func setCleanupPolicyToResourceData(cleanupPolicy *cleanuppolicies.CleanupPolicy, d *schema.ResourceData) error {
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

func resourceCleanupPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy := getCleanupPolicyFromResourceData(d)

	if err := client.CleanupPolicy.Create(&cleanupPolicy); err != nil {
		return err
	}

	d.SetId(cleanupPolicy.Name)

	return resourceCleanupPoliciesRead(d, m)
}

func resourceCleanupPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy, err := client.CleanupPolicy.Get(d.Id())
	if err != nil {
		return err
	}

	if cleanupPolicy == nil {
		d.SetId("")
		return nil
	}

	return setCleanupPolicyToResourceData(cleanupPolicy, d)
}

func resourceCleanupPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy := getCleanupPolicyFromResourceData(d)
	if err := client.CleanupPolicy.Update(&cleanupPolicy); err != nil {
		return err
	}

	return resourceCleanupPoliciesRead(d, m)
}

func resourceCleanupPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.CleanupPolicy.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceCleanupPoliciesExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	cleanupPolicy, err := client.CleanupPolicy.Get(d.Id())
	return cleanupPolicy != nil, err
}
