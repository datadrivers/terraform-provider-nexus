package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryNpmHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted Npm repository.",

		Create: resourceNpmHostedRepositoryCreate,
		Delete: resourceNpmHostedRepositoryDelete,
		Exists: resourceNpmHostedRepositoryExists,
		Read:   resourceNpmHostedRepositoryRead,
		Update: resourceNpmHostedRepositoryUpdate,
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
		},
	}
}

func getNpmHostedRepositoryFromResourceData(resourceData *schema.ResourceData) repository.NpmHostedRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))

	repo := repository.NpmHostedRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.HostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 &writePolicy,
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

func setNpmHostedRepositoryToResourceData(repo *repository.NpmHostedRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

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

func resourceNpmHostedRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getNpmHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Npm.Hosted.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceNpmHostedRepositoryRead(resourceData, m)
}

func resourceNpmHostedRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Npm.Hosted.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setNpmHostedRepositoryToResourceData(repo, resourceData)
}

func resourceNpmHostedRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getNpmHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Npm.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceNpmHostedRepositoryRead(resourceData, m)
}

func resourceNpmHostedRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Npm.Hosted.Delete(resourceData.Id())
}

func resourceNpmHostedRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Npm.Hosted.Get(resourceData.Id())
	return repo != nil, err
}
