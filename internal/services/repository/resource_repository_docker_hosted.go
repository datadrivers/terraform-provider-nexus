package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryDockerHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted docker repository.",

		Create: resourceDockerHostedRepositoryCreate,
		Delete: resourceDockerHostedRepositoryDelete,
		Exists: resourceDockerHostedRepositoryExists,
		Read:   resourceDockerHostedRepositoryRead,
		Update: resourceDockerHostedRepositoryUpdate,
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
			"storage":   repositorySchema.ResourceDockerHostedStorage,
			// Docker hosted schemas
			"docker": repositorySchema.ResourceDocker,
		},
	}
}

func getDockerHostedRepositoryFromResourceData(resourceData *schema.ResourceData) repository.DockerHostedRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	dockerConfig := resourceData.Get("docker").([]interface{})[0].(map[string]interface{})

	repo := repository.DockerHostedRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.DockerHostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 repository.StorageWritePolicy(storageConfig["write_policy"].(string)),
		},
		Docker: repository.Docker{
			ForceBasicAuth: dockerConfig["force_basic_auth"].(bool),
			V1Enabled:      dockerConfig["v1_enabled"].(bool),
		},
	}

	if latestPolicy, ok := storageConfig["latest_policy"]; ok {
		repo.Storage.LatestPolicy = tools.GetBoolPointer(latestPolicy.(bool))
	}

	if httpPort, ok := dockerConfig["http_port"]; ok {
		if httpPort.(int) > 0 {
			repo.Docker.HTTPPort = tools.GetIntPointer(httpPort.(int))
		}
	}

	if httpsPort, ok := dockerConfig["https_port"]; ok {
		if httpsPort.(int) > 0 {
			repo.Docker.HTTPSPort = tools.GetIntPointer(httpsPort.(int))
		}
	}

	if subdomain, ok := dockerConfig["subdomain"]; ok {
		if subdomain.(string) != "" {
			repo.Docker.Subdomain = tools.GetStringPointer(subdomain.(string))
		}
	}

	if pathEnabled, ok := dockerConfig["path_based_routing"]; ok {
		repo.Docker.PathEnabled = tools.GetBoolPointer(pathEnabled.(bool))
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

func setDockerHostedRepositoryToResourceData(repo *repository.DockerHostedRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if err := resourceData.Set("docker", flattenDocker(&repo.Docker)); err != nil {
		return err
	}

	if err := resourceData.Set("storage", flattenDockerHostedStorage(&repo.Storage)); err != nil {
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

func resourceDockerHostedRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getDockerHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Docker.Hosted.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceDockerHostedRepositoryRead(resourceData, m)
}

func resourceDockerHostedRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Docker.Hosted.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setDockerHostedRepositoryToResourceData(repo, resourceData)
}

func resourceDockerHostedRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getDockerHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Docker.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceDockerHostedRepositoryRead(resourceData, m)
}

func resourceDockerHostedRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Docker.Hosted.Delete(resourceData.Id())
}

func resourceDockerHostedRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Docker.Hosted.Get(resourceData.Id())
	return repo != nil, err
}
