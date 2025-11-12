package repository

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryNugetGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group nuget repository.",

		Create: resourceNugetGroupRepositoryCreate,
		Delete: resourceNugetGroupRepositoryDelete,
		Exists: resourceNugetGroupRepositoryExists,
		Read:   resourceNugetGroupRepositoryRead,
		Update: resourceNugetGroupRepositoryUpdate,
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

func getNugetGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.NugetGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.NugetGroupRepository{
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

func setNugetGroupRepositoryToResourceData(repo *repository.NugetGroupRepository, resourceData *schema.ResourceData) error {
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

func resourceNugetGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getNugetGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Nuget.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceNugetGroupRepositoryRead(resourceData, m)
}

func resourceNugetGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Nuget.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setNugetGroupRepositoryToResourceData(repo, resourceData)
}

func resourceNugetGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getNugetGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Nuget.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceNugetGroupRepositoryRead(resourceData, m)
}

func resourceNugetGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Nuget.Group.Delete(resourceData.Id())
}

func resourceNugetGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Nuget.Group.Get(resourceData.Id())
	return repo != nil, err
}
