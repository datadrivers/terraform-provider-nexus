package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryBowerGroup() *schema.Resource {
	return &schema.Resource{
		Description: `!> This resource is deprecated and will be removed in the next major release of this provider. Bower repositories were removed in Nexus 3.71.0.

Use this resource to create a group bower repository.`,
		DeprecationMessage: "This resource is deprecated and will be removed in the next major release of this provider. Bower repositories were removed in Nexus 3.71.0.",

		Create: resourceBowerGroupRepositoryCreate,
		Delete: resourceBowerGroupRepositoryDelete,
		Exists: resourceBowerGroupRepositoryExists,
		Read:   resourceBowerGroupRepositoryRead,
		Update: resourceBowerGroupRepositoryUpdate,
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

func getBowerGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.BowerGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.BowerGroupRepository{
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

func setBowerGroupRepositoryToResourceData(repo *repository.BowerGroupRepository, resourceData *schema.ResourceData) error {
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

func resourceBowerGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getBowerGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Bower.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceBowerGroupRepositoryRead(resourceData, m)
}

func resourceBowerGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Bower.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setBowerGroupRepositoryToResourceData(repo, resourceData)
}

func resourceBowerGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getBowerGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Bower.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceBowerGroupRepositoryRead(resourceData, m)
}

func resourceBowerGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Bower.Group.Delete(resourceData.Id())
}

func resourceBowerGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Bower.Group.Get(resourceData.Id())
	return repo != nil, err
}
