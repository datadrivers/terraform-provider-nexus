package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted yum repository.",

		Create: resourceYumHostedRepositoryCreate,
		Read:   resourceYumHostedRepositoryRead,
		Update: resourceYumHostedRepositoryUpdate,
		Delete: resourceYumHostedRepositoryDelete,
		Exists: resourceYumHostedRepositoryExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cleanup": getResourceCleanupSchema(),
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
			"storage": getResourceHostedStorageSchema(),
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
					PolicyNames: tools.InterfaceSliceToStringSlice(policy_names.(*schema.Set).List()),
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

	if err := d.Set("storage", flattenHostedStorage(&repo.Storage, d)); err != nil {
		return err
	}

	if repo.Cleanup != nil {
		if err := d.Set("cleanup", flattenCleanup(repo.Cleanup)); err != nil {
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
