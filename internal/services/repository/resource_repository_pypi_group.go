package repository

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/dre2004/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryPypiGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group pypi repository.",

		Create: resourcePypiGroupRepositoryCreate,
		Delete: resourcePypiGroupRepositoryDelete,
		Exists: resourcePypiGroupRepositoryExists,
		Read:   resourcePypiGroupRepositoryRead,
		Update: resourcePypiGroupRepositoryUpdate,
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

func getPypiGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.PypiGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].(*schema.Set).List() {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.PypiGroupRepository{
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

func setPypiGroupRepositoryToResourceData(repo *repository.PypiGroupRepository, resourceData *schema.ResourceData) error {
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

func resourcePypiGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getPypiGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Pypi.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourcePypiGroupRepositoryRead(resourceData, m)
}

func resourcePypiGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Pypi.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setPypiGroupRepositoryToResourceData(repo, resourceData)
}

func resourcePypiGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getPypiGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Pypi.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourcePypiGroupRepositoryRead(resourceData, m)
}

func resourcePypiGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Pypi.Group.Delete(resourceData.Id())
}

func resourcePypiGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Pypi.Group.Get(resourceData.Id())
	return repo != nil, err
}
