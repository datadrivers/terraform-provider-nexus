package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryMavenHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted maven repository.",

		Create: resourceMavenHostedRepositoryCreate,
		Delete: resourceMavenHostedRepositoryDelete,
		Exists: resourceMavenHostedRepositoryExists,
		Read:   resourceMavenHostedRepositoryRead,
		Update: resourceMavenHostedRepositoryUpdate,
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
			// Maven hosted schemas
			"maven": repositorySchema.ResourceMaven,
		},
	}
}

func getMavenHostedRepositoryFromResourceData(resourceData *schema.ResourceData) repository.MavenHostedRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))
	mavenConfig := resourceData.Get("maven").([]interface{})[0].(map[string]interface{})

	repo := repository.MavenHostedRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.HostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 &writePolicy,
		},
		Maven: repository.Maven{
			VersionPolicy: repository.MavenVersionPolicy(mavenConfig["version_policy"].(string)),
			LayoutPolicy:  repository.MavenLayoutPolicy(mavenConfig["layout_policy"].(string)),
		},
	}

	if mavenConfig["content_disposition"] != "" {
		contentDisposition := repository.MavenContentDisposition(mavenConfig["content_disposition"].(string))
		repo.Maven.ContentDisposition = &contentDisposition
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

func setMavenHostedRepositoryToResourceData(repo *repository.MavenHostedRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if err := resourceData.Set("storage", flattenHostedStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("maven", flattenMaven(&repo.Maven)); err != nil {
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

func resourceMavenHostedRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getMavenHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Maven.Hosted.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceMavenHostedRepositoryRead(resourceData, m)
}

func resourceMavenHostedRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Maven.Hosted.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setMavenHostedRepositoryToResourceData(repo, resourceData)
}

func resourceMavenHostedRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getMavenHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Maven.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceMavenHostedRepositoryRead(resourceData, m)
}

func resourceMavenHostedRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Maven.Hosted.Delete(resourceData.Id())
}

func resourceMavenHostedRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Maven.Hosted.Get(resourceData.Id())
	return repo != nil, err
}
