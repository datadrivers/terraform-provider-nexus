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

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Create: resourceYumHostedRepositoryCreate,
		Read:   resourceYumHostedRepositoryRead,
		Update: resourceYumHostedRepositoryUpdate,
		Delete: resourceYumHostedRepositoryDelete,
		Exists: resourceYumHostedRepositoryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cleanup": {
				DefaultFunc: repositoryCleanupDefault,
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
				Required:    true,
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
		},
	}
}

func getYumHostedRepositoryFromResourceData(d *schema.ResourceData) repository.YumHostedRepository {
	storageConfig := d.Get("storage").([]interface{})[0].(map[string]interface{})
	writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))
	deployPolicy := repository.YumDeployPolicy(d.Get("deploy_policy").(string))

	repo := repository.YumHostedRepository{
		Name:   d.Get("name").(string),
		Online: d.Get("online").(bool),
		Storage: repository.HostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 &writePolicy,
		},
		Yum: repository.Yum{
			RepodataDepth: d.Get("repodata_depth").(int),
			DeployPolicy:  &deployPolicy,
		},
	}

	cleanupList := d.Get("cleanup").([]interface{})
	if len(cleanupList) > 0 && cleanupList[0] != nil {
		cleanupConfig := cleanupList[0].(map[string]interface{})
		if len(cleanupConfig) > 0 {
			policy_names, ok := cleanupConfig["policy_names"]
			if ok {
				repo.Cleanup = &repository.Cleanup{
					PolicyNames: interfaceSliceToStringSlice(policy_names.(*schema.Set).List()),
				}
			}
		}
	}

	return repo
}

func setYumHostedRepositoryToResourceData(repo *repository.YumHostedRepository, d *schema.ResourceData) error {
	d.SetId(repo.Name)
	d.Set("name", repo.Name)
	d.Set("online", repo.Online)
	d.Set("repodata_depth", repo.Yum.RepodataDepth)
	d.Set("deploy_policy", repo.Yum.DeployPolicy)

	if err := d.Set("storage", flattenRepositoryHostedStorage(&repo.Storage, d)); err != nil {
		return err
	}

	if repo.Cleanup != nil {
		if err := d.Set("cleanup", flattenRepositoryCleanup(repo.Cleanup)); err != nil {
			return err
		}
	}

	return nil
}

func resourceYumHostedRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getYumHostedRepositoryFromResourceData(d)

	if err := client.Repository.Yum.Hosted.Create(repo); err != nil {
		return err
	}
	d.SetId(repo.Name)

	return resourceYumHostedRepositoryRead(d, m)
}

func resourceYumHostedRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Hosted.Get(d.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		d.SetId("")
		return nil
	}

	return setYumHostedRepositoryToResourceData(repo, d)
}

func resourceYumHostedRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := d.Id()
	repo := getYumHostedRepositoryFromResourceData(d)

	if err := client.Repository.Yum.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceYumHostedRepositoryRead(d, m)
}

func resourceYumHostedRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Yum.Hosted.Delete(d.Id())
}

func resourceYumHostedRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Hosted.Get(d.Id())
	return repo != nil, err
}
