package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryCargoHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted Cargo repository.",

		Create: resourceCargoHostedRepositoryCreate,
		Delete: resourceCargoHostedRepositoryDelete,
		Exists: resourceCargoHostedRepositoryExists,
		Read:   resourceCargoHostedRepositoryRead,
		Update: resourceCargoHostedRepositoryUpdate,
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

func getCargoHostedRepositoryFromResourceData(resourceData *schema.ResourceData) repository.CargoHostedRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))

	repo := repository.CargoHostedRepository{
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
			policyNames, ok := cleanupConfig["policy_names"]
			if ok {
				repo.Cleanup = &repository.Cleanup{
					PolicyNames: tools.InterfaceSliceToStringSlice(policyNames.(*schema.Set).List()),
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

func setCargoHostedRepositoryToResourceData(repo *repository.CargoHostedRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	if err := resourceData.Set("name", repo.Name); err != nil {
		return err
	}
	if err := resourceData.Set("online", repo.Online); err != nil {
		return err
	}

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

func resourceCargoHostedRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getCargoHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Cargo.Hosted.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceCargoHostedRepositoryRead(resourceData, m)
}

func resourceCargoHostedRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Cargo.Hosted.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setCargoHostedRepositoryToResourceData(repo, resourceData)
}

func resourceCargoHostedRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getCargoHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Cargo.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceCargoHostedRepositoryRead(resourceData, m)
}

func resourceCargoHostedRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Cargo.Hosted.Delete(resourceData.Id())
}

func resourceCargoHostedRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Cargo.Hosted.Get(resourceData.Id())
	return repo != nil, err
}
