package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryGoGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group go repository.",

		Create: resourceGoGroupRepositoryCreate,
		Delete: resourceGoGroupRepositoryDelete,
		Exists: resourceGoGroupRepositoryExists,
		Read:   resourceGoGroupRepositoryRead,
		Update: resourceGoGroupRepositoryUpdate,
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

func getGoGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.GoGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.GoGroupRepository{
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

func setGoGroupRepositoryToResourceData(repo *repository.GoGroupRepository, resourceData *schema.ResourceData) error {
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

func resourceGoGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getGoGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Go.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceGoGroupRepositoryRead(resourceData, m)
}

func resourceGoGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Go.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setGoGroupRepositoryToResourceData(repo, resourceData)
}

func resourceGoGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getGoGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Go.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceGoGroupRepositoryRead(resourceData, m)
}

func resourceGoGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Go.Group.Delete(resourceData.Id())
}

func resourceGoGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Go.Group.Get(resourceData.Id())
	return repo != nil, err
}
