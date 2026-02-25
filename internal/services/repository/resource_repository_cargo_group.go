package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryCargoGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group cargo repository.",

		Create: resourceCargoGroupRepositoryCreate,
		Delete: resourceCargoGroupRepositoryDelete,
		Exists: resourceCargoGroupRepositoryExists,
		Read:   resourceCargoGroupRepositoryRead,
		Update: resourceCargoGroupRepositoryUpdate,
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

func getCargoGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.CargoGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.CargoGroupRepository{
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

func setCargoGroupRepositoryToResourceData(repo *repository.CargoGroupRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	if err := resourceData.Set("name", repo.Name); err != nil {
		return err
	}
	if err := resourceData.Set("online", repo.Online); err != nil {
		return err
	}

	if err := resourceData.Set("storage", flattenStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("group", flattenGroup(&repo.Group)); err != nil {
		return err
	}

	return nil
}

func resourceCargoGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getCargoGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Cargo.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceCargoGroupRepositoryRead(resourceData, m)
}

func resourceCargoGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Cargo.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setCargoGroupRepositoryToResourceData(repo, resourceData)
}

func resourceCargoGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getCargoGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Cargo.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceCargoGroupRepositoryRead(resourceData, m)
}

func resourceCargoGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Cargo.Group.Delete(resourceData.Id())
}

func resourceCargoGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Cargo.Group.Get(resourceData.Id())
	return repo != nil, err
}
