package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryRubygemsGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group rubygems repository.",

		Create: resourceRubygemsGroupRepositoryCreate,
		Delete: resourceRubygemsGroupRepositoryDelete,
		Exists: resourceRubygemsGroupRepositoryExists,
		Read:   resourceRubygemsGroupRepositoryRead,
		Update: resourceRubygemsGroupRepositoryUpdate,
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

func getRubygemsGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.RubyGemsGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.RubyGemsGroupRepository{
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

func setRubygemsGroupRepositoryToResourceData(repo *repository.RubyGemsGroupRepository, resourceData *schema.ResourceData) error {
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

func resourceRubygemsGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getRubygemsGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.RubyGems.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceRubygemsGroupRepositoryRead(resourceData, m)
}

func resourceRubygemsGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.RubyGems.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setRubygemsGroupRepositoryToResourceData(repo, resourceData)
}

func resourceRubygemsGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getRubygemsGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.RubyGems.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceRubygemsGroupRepositoryRead(resourceData, m)
}

func resourceRubygemsGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.RubyGems.Group.Delete(resourceData.Id())
}

func resourceRubygemsGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.RubyGems.Group.Get(resourceData.Id())
	return repo != nil, err
}
