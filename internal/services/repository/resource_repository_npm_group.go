package repository

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/dre2004/terraform-provider-nexus/internal/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryNpmGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group npm repository.",

		Create: resourceNpmGroupRepositoryCreate,
		Delete: resourceNpmGroupRepositoryDelete,
		Exists: resourceNpmGroupRepositoryExists,
		Read:   resourceNpmGroupRepositoryRead,
		Update: resourceNpmGroupRepositoryUpdate,
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
		},
	}
}

func getNpmGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.NpmGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].(*schema.Set).List() {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.NpmGroupRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.Storage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		},
		Group: repository.GroupDeploy{
			MemberNames: groupMemberNames,
		},
	}

	if groupConfig["writable_member"].(string) != "" {
		repo.Group.WritableMember = tools.GetStringPointer(groupConfig["writable_member"].(string))
	}

	return repo
}

func setNpmGroupRepositoryToResourceData(repo *repository.NpmGroupRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if err := resourceData.Set("storage", flattenStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("group", flattenGroupDeploy(&repo.Group)); err != nil {
		return err
	}

	return nil
}

func resourceNpmGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getNpmGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Npm.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceNpmGroupRepositoryRead(resourceData, m)
}

func resourceNpmGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Npm.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setNpmGroupRepositoryToResourceData(repo, resourceData)
}

func resourceNpmGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getNpmGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Npm.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceNpmGroupRepositoryRead(resourceData, m)
}

func resourceNpmGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Npm.Group.Delete(resourceData.Id())
}

func resourceNpmGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Npm.Group.Get(resourceData.Id())
	return repo != nil, err
}
