package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryDockerGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group docker repository.",

		Create: resourceDockerGroupRepositoryCreate,
		Delete: resourceDockerGroupRepositoryDelete,
		Exists: resourceDockerGroupRepositoryExists,
		Read:   resourceDockerGroupRepositoryRead,
		Update: resourceDockerGroupRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.ResourceID,
			"name":   repositorySchema.ResourceName,
			"online": repositorySchema.ResourceOnline,
			// Group schemas
			"group":   repositorySchema.ResourceGroupDeploy,
			"storage": repositorySchema.ResourceStorage,
			// Docker group schemas
			"docker": repositorySchema.ResourceDocker,
		},
	}
}

func getDockerGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.DockerGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	dockerConfig := resourceData.Get("docker").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].(*schema.Set).List() {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.DockerGroupRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.Storage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		},
		Group: repository.GroupDeploy{
			MemberNames: groupMemberNames,
		},
		Docker: repository.Docker{
			ForceBasicAuth: dockerConfig["force_basic_auth"].(bool),
			V1Enabled:      dockerConfig["v1_enabled"].(bool),
		},
	}

	if groupConfig["writable_member"].(string) != "" {
		repo.Group.WritableMember = tools.GetStringPointer(groupConfig["writable_member"].(string))
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

	return repo
}

func setDockerGroupRepositoryToResourceData(repo *repository.DockerGroupRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if err := resourceData.Set("docker", flattenDocker(&repo.Docker)); err != nil {
		return err
	}

	if err := resourceData.Set("storage", flattenStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("group", flattenGroupDeploy(&repo.Group)); err != nil {
		return err
	}

	return nil
}

func resourceDockerGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getDockerGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Docker.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceDockerGroupRepositoryRead(resourceData, m)
}

func resourceDockerGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Docker.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setDockerGroupRepositoryToResourceData(repo, resourceData)
}

func resourceDockerGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getDockerGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Docker.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceDockerGroupRepositoryRead(resourceData, m)
}

func resourceDockerGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Docker.Group.Delete(resourceData.Id())
}

func resourceDockerGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Docker.Group.Get(resourceData.Id())
	return repo != nil, err
}
