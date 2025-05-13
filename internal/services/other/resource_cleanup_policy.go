package other

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	RepositoryFormats = []string{
		"*",
		"apt",
		"cargo",
		"cocoapods",
		"composer",
		"conan",
		"conda",
		"docker",
		"gitlfs",
		"go",
		"helm",
		"huggingface",
		"maven2",
		"npm",
		"nuget",
		"p2",
		"pypi",
		"r",
		"raw",
		"rubygems",
		"yum",
	}
	ReleaseTypes = []string{
		"RELEASES",
		"PRERELEASES",
		"RELEASES_AND_PRERELEASES",
	}
)

func ResourceCleanupPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create and execute a custom Cleanup Policy.",

		Create: resourceCleanupPolicyCreate,
		Read:   resourceCleanupPolicyRead,
		Update: resourceCleanupPolicyUpdate,
		Delete: resourceCleanupPolicyDelete,
		Exists: resourceCleanupPolicyExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"name": {
				Description:  "Use a unique name for the cleanup policy",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"format": {
				Description:  "The format that this cleanup policy can be applied to",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(RepositoryFormats, false),
			},
			"retain": {
				Description: "number of versions to keep. Require: Nexus Pro",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"criteria": {
				Description: "Remove all components that match all selected criteria.",
				Required:    true,
				MaxItems:    1,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_blob_updated": {
							Description: "The age of the component in days.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"last_downloaded": {
							Description: "The last time the component had been downloaded in days.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"release_type": {
							Description:  "Is one of: RELEASES_AND_PRERELEASES, PRERELEASES, RELEASES], Only maven2, npm and yum repositories support this field.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(ReleaseTypes, false),
						},
						"asset_regex": {
							Description: "A regex string to filter for specific asset paths. Not for gitlfs or *",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func getCleanupPolicyFromResourceData(d *schema.ResourceData) *nexusSchema.CleanupPolicy {
	criteria := d.Get("criteria").([]any)[0].(map[string]any)

	cp := &nexusSchema.CleanupPolicy{
		Name:   d.Get("name").(string),
		Format: nexusSchema.RepositoryFormat(d.Get("format").(string)),
		Retain: d.Get("retain").(int),
	}

	if c, ok := criteria["last_blob_updated"]; ok {
		cp.CriteriaLastBlobUpdated = tools.GetIntPointer(c.(int))
	}
	if c, ok := criteria["last_downloaded"]; ok {
		cp.CriteriaLastDownloaded = tools.GetIntPointer(c.(int))
	}
	if c, ok := criteria["release_type"]; ok {
		cr := nexusSchema.CriteriaReleaseType(c.(string))
		cp.CriteriaReleaseType = &cr
	}
	if c, ok := criteria["asset_regex"]; ok {
		cp.CriteriaAssetRegex = tools.GetStringPointer(c.(string))
	}
	return cp
}

func setCleanupPolicyToResourceData(cp *nexusSchema.CleanupPolicy, d *schema.ResourceData) error {
	d.SetId(cp.Name)
	if err := d.Set("name", cp.Name); err != nil {
		return err
	}
	if err := d.Set("format", cp.Format); err != nil {
		return err
	}
	if err := d.Set("retain", cp.Retain); err != nil {
		return err
	}

	criteria := make(map[string]any)
	if c := cp.CriteriaLastBlobUpdated; c != nil {
		criteria["last_blob_updated"] = *c
	}
	if c := cp.CriteriaLastDownloaded; c != nil {
		criteria["last_downloaded"] = *c
	}
	if c := cp.CriteriaReleaseType; c != nil {
		criteria["release_type"] = string(*c)
	}
	if c := cp.CriteriaAssetRegex; c != nil {
		criteria["asset_regex"] = *c
	}

	if len(criteria) > 0 {
		if err := d.Set("criteria", []any{criteria}); err != nil {
			return err
		}
	}

	return nil
}

func resourceCleanupPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	cp := getCleanupPolicyFromResourceData(d)

	if err := client.CleanupPolicy.Create(cp); err != nil {
		return err
	}

	d.SetId(cp.Name)
	return resourceCleanupPolicyRead(d, m)
}

func resourceCleanupPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cp, err := client.CleanupPolicy.Get(d.Id())
	if err != nil {
		return err
	}

	if cp == nil {
		d.SetId("")
		return nil
	}

	return setCleanupPolicyToResourceData(cp, d)
}

func resourceCleanupPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cleanupPolicy := getCleanupPolicyFromResourceData(d)
	if err := client.CleanupPolicy.Update(cleanupPolicy); err != nil {
		return err
	}

	return resourceCleanupPolicyRead(d, m)
}

func resourceCleanupPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	return client.CleanupPolicy.Delete(d.Id())
}

func resourceCleanupPolicyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	cleanupPolicy, err := client.CleanupPolicy.Get(d.Id())
	return cleanupPolicy != nil, err
}
