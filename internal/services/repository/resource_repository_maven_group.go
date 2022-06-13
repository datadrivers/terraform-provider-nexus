package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryMavenGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group maven repository.",

		Create: resourceMavenGroupRepositoryCreate,
		Delete: resourceMavenGroupRepositoryDelete,
		Exists: resourceMavenGroupRepositoryExists,
		Read:   resourceMavenGroupRepositoryRead,
		Update: resourceMavenGroupRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.ResourceID,
			"name":   repositorySchema.ResourceName,
			"online": repositorySchema.ResourceOnline,
			// Group schemas
			"group":   repositorySchema.ResourceGroup,
			"storage": repositorySchema.ResourceStorage,
		},
	}
}

func getMavenGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.MavenGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].(*schema.Set).List() {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.MavenGroupRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.Storage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		},
		Group: repository.Group{
			MemberNames: groupMemberNames,
		},
	}

	return repo
}

func setMavenGroupRepositoryToResourceData(repo *repository.MavenGroupRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if err := resourceData.Set("storage", flattenStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("group", flattenGroup(&repo.Group)); err != nil {
		return err
	}

	return nil
}

func resourceMavenGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getMavenGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Maven.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceMavenGroupRepositoryRead(resourceData, m)
}

func resourceMavenGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Maven.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setMavenGroupRepositoryToResourceData(repo, resourceData)
}

func resourceMavenGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getMavenGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Maven.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceMavenGroupRepositoryRead(resourceData, m)
}

func resourceMavenGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Maven.Group.Delete(resourceData.Id())
}

func resourceMavenGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Maven.Group.Get(resourceData.Id())
	return repo != nil, err
}
