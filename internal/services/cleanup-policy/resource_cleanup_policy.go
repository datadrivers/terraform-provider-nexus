package cleanup_policy

import (
	"fmt"
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCleanupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCleanupPolicyCreate,
		Read:   resourceCleanupPolicyRead,
		Update: resourceCleanupPolicyUpdate,
		Delete: resourceCleanupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"criteria_last_blob_updated": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"criteria_last_downloaded": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"criteria_release_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RELEASES_AND_PRERELEASES",
					"PRERELEASES",
					"RELEASES",
				}, false),
			},
			"criteria_asset_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retain": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"format": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"*", "apt", "bower", "cocoapods", "conan", "conda", "docker", "gitlfs", "go", "helm", "maven2", "npm", "nuget", "p2", "pypi", "r", "raw", "rubygems", "yum",
				}, false),
			},
		},
	}
}

func resourceCleanupPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	policy := nexusSchema.CleanupPolicy{
		Name:                    d.Get("name").(string),
		Notes:                   getStringPointer(d, "notes"),
		CriteriaLastBlobUpdated: getIntPointer(d, "criteria_last_blob_updated"),
		CriteriaLastDownloaded:  getIntPointer(d, "criteria_last_downloaded"),
		CriteriaReleaseType:     getCriteriaReleaseTypePointer(d, "criteria_release_type"),
		CriteriaAssetRegex:      getStringPointer(d, "criteria_asset_regex"),
		Retain:                  getInt(d, "retain"),
		Format:                  nexusSchema.RepositoryFormat(d.Get("format").(string)),
	}

	err := client.CleanupPolicy.Create(&policy)
	if err != nil {
		return err
	}
	d.SetId(policy.Name)
	return resourceCleanupPolicyRead(d, m)
}

func resourceCleanupPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	policy, err := client.CleanupPolicy.Get(d.Id())
	if err != nil {
		return err
	}

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

func resourceCleanupPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") {
		return fmt.Errorf("updating the name of a cleanup policy is not allowed")
	}

	client := m.(*nexus.NexusClient)
	policy := nexusSchema.CleanupPolicy{
		Name:                    d.Get("name").(string),
		Notes:                   getStringPointer(d, "notes"),
		CriteriaLastBlobUpdated: getIntPointer(d, "criteria_last_blob_updated"),
		CriteriaLastDownloaded:  getIntPointer(d, "criteria_last_downloaded"),
		CriteriaReleaseType:     getCriteriaReleaseTypePointer(d, "criteria_release_type"),
		CriteriaAssetRegex:      getStringPointer(d, "criteria_asset_regex"),
		Retain:                  getInt(d, "retain"),
		Format:                  nexusSchema.RepositoryFormat(d.Get("format").(string)),
	}

	err := client.CleanupPolicy.Update(&policy)
	if err != nil {
		return err
	}
	return resourceCleanupPolicyRead(d, m)
}

func resourceCleanupPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.CleanupPolicy.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func getIntPointer(d *schema.ResourceData, key string) *int {
	if v, ok := d.GetOk(key); ok {
		value := v.(int)
		return &value
	}
	return nil
}

func getStringPointer(d *schema.ResourceData, key string) *string {
	if v, ok := d.GetOk(key); ok {
		value := v.(string)
		return &value
	}
	return nil
}

func getInt(d *schema.ResourceData, key string) int {
	if v, ok := d.GetOk(key); ok {
		return v.(int)
	}
	return 0
}

func getCriteriaReleaseTypePointer(d *schema.ResourceData, key string) *nexusSchema.CriteriaReleaseType {
	if v, ok := d.GetOk(key); ok {
		value := nexusSchema.CriteriaReleaseType(v.(string))
		return &value
	}
	return nil
}
