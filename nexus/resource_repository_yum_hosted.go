/*
Use this resource to create a hosted yum repository.

Example Usage

```hcl

resource "nexus_repository_yum_hosted" "yum" {
  name = "yummy"
}

resource "nexus_repository_yum_hosted" "yum1" {
  deploy_policy  = "STRICT"
  name = "yummy1"
  online = true
  repodata_depth = 4

  cleanup {
    policy_names = ["policy"]
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW"
  }

}
```
*/
package nexus

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Create: resourceYumRepositoryCreate,
		Read:   resourceYumRepositoryRead,
		Update: resourceYumRepositoryUpdate,
		Delete: resourceYumRepositoryDelete,
		Exists: resourceYumRepositoryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cleanup": {
				DefaultFunc: RepositoryCleanupDefault,
				Description: "Cleanup policies",
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Description: "List of policy names",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Set: func(v interface{}) int {
								return schema.HashString(strings.ToLower(v.(string)))
							},
							Type: schema.TypeSet,
						},
					},
				},
			},
			"deploy_policy": {
				Default:     "STRICT",
				Description: "Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Default:     true,
				Description: "Whether this repository accepts incoming requests",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"repodata_depth": {
				Default:     0,
				Description: "Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			"storage": {
				DefaultFunc: repositoryStorageDefault,
				Description: "The storage configuration of the repository",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_store_name": {
							Default:     "default",
							Description: "Blob store used to store repository contents",
							Optional:    true,
							Set: func(v interface{}) int {
								return schema.HashString(strings.ToLower(v.(string)))
							},
							Type: schema.TypeString,
						},
						"strict_content_type_validation": {
							Default:     true,
							Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"write_policy": {
							Description: "Controls if deployments of and updates to assets are allowed",
							Default:     "ALLOW",
							Optional:    true,
							Type:        schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"ALLOW",
								"ALLOW_ONCE",
								"DENY",
							}, false),
						},
					},
				},
			},
			"type": {
				Description: "Repository type",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func getYumRepositoryFromResourceData(d *schema.ResourceData) nexus.Repository {
	repo := nexus.Repository{
		Format: "yum",
		Name:   d.Get("name").(string),
		Online: d.Get("online").(bool),
		Type:   "hosted",
	}

	repo.RepositoryCleanup = &nexus.RepositoryCleanup{
		PolicyNames: []string{},
	}

	cleanupList := d.Get("cleanup").([]interface{})
	if len(cleanupList) > 0 && cleanupList[0] != nil {
		cleanupConfig := cleanupList[0].(map[string]interface{})
		if len(cleanupConfig) > 0 {
			policy_names, ok := cleanupConfig["policy_names"]
			if ok {
				repo.RepositoryCleanup = &nexus.RepositoryCleanup{
					PolicyNames: interfaceSliceToStringSlice(policy_names.(*schema.Set).List()),
				}
			}
		}
	}

	storageList := d.Get("storage").([]interface{})
	if len(storageList) > 0 {
		storageConfig := storageList[0].(map[string]interface{})

		writePolicy := storageConfig["write_policy"].(string)

		repo.RepositoryStorage = &nexus.RepositoryStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 &writePolicy,
		}
	} else {
		writePolicy := "ALLOW"
		repo.RepositoryStorage = &nexus.RepositoryStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		}

	}

	repo.RepositoryYum = &nexus.RepositoryYum{
		RepodataDepth: d.Get("repodata_depth").(int),
		DeployPolicy:  d.Get("deploy_policy").(string),
	}

	return repo
}

func setYumRepositoryToResourceData(repo *nexus.Repository, d *schema.ResourceData) error {
	d.SetId(repo.Name)
	d.Set("name", repo.Name)
	d.Set("online", repo.Online)
	d.Set("type", "hosted")

	if repo.RepositoryCleanup != nil {
		if err := d.Set("cleanup", flattenRepositoryCleanup(repo.RepositoryCleanup)); err != nil {
			return err
		}
	}

	d.Set("repodata_depth", repo.RepositoryYum.RepodataDepth)
	d.Set("deploy_policy", repo.RepositoryYum.DeployPolicy)

	if repo.RepositoryStorage != nil {
		if err := d.Set("storage", flattenRepositoryStorage(repo.RepositoryStorage, d)); err != nil {
			return err
		}
	}

	return nil
}

func resourceYumRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	repo := getYumRepositoryFromResourceData(d)

	if err := client.RepositoryCreate(repo); err != nil {
		return err
	}

	if err := setYumRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceYumRepositoryRead(d, m)
}

func resourceYumRepositoryRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	repo, err := nexusClient.RepositoryRead(d.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		d.SetId("")
		return nil
	}

	return setYumRepositoryToResourceData(repo, d)
}

func resourceYumRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	repoName := d.Id()
	repo := getYumRepositoryFromResourceData(d)

	if err := client.RepositoryUpdate(repoName, repo); err != nil {
		return err
	}

	if err := setYumRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceYumRepositoryRead(d, m)
}

func resourceYumRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	return nexusClient.RepositoryDelete(d.Id())
}

func resourceYumRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	repo, err := nexusClient.RepositoryRead(d.Id())
	return repo != nil, err
}
