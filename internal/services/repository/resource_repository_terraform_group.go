package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryTerraformGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group terraform repository.",

		Create: resourceTerraformGroupRepositoryCreate,
		Delete: resourceTerraformGroupRepositoryDelete,
		Exists: resourceTerraformGroupRepositoryExists,
		Read:   resourceTerraformGroupRepositoryRead,
		Update: resourceTerraformGroupRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id":      common.ResourceID,
			"name":    repositorySchema.ResourceName,
			"online":  repositorySchema.ResourceOnline,
			"group":   repositorySchema.ResourceGroup,
			"storage": repositorySchema.ResourceStorage,
		},
	}
}

func getTerraformGroupRepositoryFromResourceData(resourceData *schema.ResourceData) TerraformGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})

	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	return TerraformGroupRepository{
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
}

func setTerraformGroupRepositoryToResourceData(repo *TerraformGroupRepository, resourceData *schema.ResourceData) error {
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

func resourceTerraformGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getTerraformGroupRepositoryFromResourceData(resourceData)

	if err := terraformGroupService(client).Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceTerraformGroupRepositoryRead(resourceData, m)
}

func resourceTerraformGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := terraformGroupService(client).Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setTerraformGroupRepositoryToResourceData(repo, resourceData)
}

func resourceTerraformGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getTerraformGroupRepositoryFromResourceData(resourceData)

	if err := terraformGroupService(client).Update(repoName, repo); err != nil {
		return err
	}

	return resourceTerraformGroupRepositoryRead(resourceData, m)
}

func resourceTerraformGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return terraformGroupService(client).Delete(resourceData.Id())
}

func resourceTerraformGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := terraformGroupService(client).Get(resourceData.Id())
	return repo != nil, err
}
