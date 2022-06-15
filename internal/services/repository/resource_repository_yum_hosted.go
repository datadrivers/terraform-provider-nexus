package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceRepositoryYumHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted yum repository.",

		Create: resourceYumHostedRepositoryCreate,
		Delete: resourceYumHostedRepositoryDelete,
		Exists: resourceYumHostedRepositoryExists,
		Read:   resourceYumHostedRepositoryRead,
		Update: resourceYumHostedRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.ResourceID,
			"name":   repositorySchema.ResourceName,
			"online": repositorySchema.ResourceOnline,
			// Hosted schemas
			"cleanup":   repositorySchema.ResourceCleanup,
			"component": repositorySchema.ResourceComponent,
			"storage":   repositorySchema.ResourceHostedStorage,
			// Yum hosted schemas
			"deploy_policy": {
				Default:      "STRICT",
				Description:  "Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(repository.YumDeployPolicyStrict), string(repository.YumDeployPolicyPermissive)}, false),
			},
			"repodata_depth": {
				Default:      0,
				Description:  "Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5",
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 5),
			},
		},
	}
}

func getYumHostedRepositoryFromResourceData(resourceData *schema.ResourceData) repository.YumHostedRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))
	deployPolicy := repository.YumDeployPolicy(resourceData.Get("deploy_policy").(string))

	repo := repository.YumHostedRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.HostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 &writePolicy,
		},
		Yum: repository.Yum{
			RepodataDepth: resourceData.Get("repodata_depth").(int),
			DeployPolicy:  &deployPolicy,
		},
	}

	cleanupList := resourceData.Get("cleanup").([]interface{})
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

	componentList := resourceData.Get("component").([]interface{})
	if len(componentList) > 0 && componentList[0] != nil {
		componentConfig := componentList[0].(map[string]interface{})
		if len(componentConfig) > 0 {
			repo.Component = &repository.Component{
				ProprietaryComponents: componentConfig["proprietary_components"].(bool),
			}
		}
	}

	return repo
}

func setYumHostedRepositoryToResourceData(repo *repository.YumHostedRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)
	resourceData.Set("repodata_depth", repo.Yum.RepodataDepth)
	resourceData.Set("deploy_policy", repo.Yum.DeployPolicy)

	if err := resourceData.Set("storage", flattenHostedStorage(&repo.Storage)); err != nil {
		return err
	}

	if repo.Cleanup != nil {
		if err := resourceData.Set("cleanup", flattenCleanup(repo.Cleanup)); err != nil {
			return err
		}
	}

	if repo.Component != nil {
		if err := resourceData.Set("component", flattenComponent(repo.Component)); err != nil {
			return err
		}
	}

	return nil
}

func resourceYumHostedRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getYumHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Yum.Hosted.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceYumHostedRepositoryRead(resourceData, m)
}

func resourceYumHostedRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Hosted.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setYumHostedRepositoryToResourceData(repo, resourceData)
}

func resourceYumHostedRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getYumHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Yum.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceYumHostedRepositoryRead(resourceData, m)
}

func resourceYumHostedRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Yum.Hosted.Delete(resourceData.Id())
}

func resourceYumHostedRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Hosted.Get(resourceData.Id())
	return repo != nil, err
}
