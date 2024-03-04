package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryYumGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a group yum repository.",

		Create: resourceYumGroupRepositoryCreate,
		Delete: resourceYumGroupRepositoryDelete,
		Exists: resourceYumGroupRepositoryExists,
		Read:   resourceYumGroupRepositoryRead,
		Update: resourceYumGroupRepositoryUpdate,
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
			// Yum group schemas
			"yum_signing": repositorySchema.ResourceYumSigning,
		},
	}
}

func getYumGroupRepositoryFromResourceData(resourceData *schema.ResourceData) repository.YumGroupRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	groupConfig := resourceData.Get("group").([]interface{})[0].(map[string]interface{})
	groupMemberNames := []string{}
	for _, name := range groupConfig["member_names"].([]interface{}) {
		groupMemberNames = append(groupMemberNames, name.(string))
	}

	repo := repository.YumGroupRepository{
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

	yumSigningList := resourceData.Get("yum_signing").([]interface{})
	if len(yumSigningList) > 0 && yumSigningList[0] != nil {
		yumSigningConfig := yumSigningList[0].(map[string]interface{})

		repo.YumSigning = &repository.YumSigning{
			Keypair: tools.GetStringPointer(yumSigningConfig["keypair"].(string)),
		}
		if yumSigningConfig["passphrase"].(string) != "" {
			repo.YumSigning.Passphrase = tools.GetStringPointer(yumSigningConfig["passphrase"].(string))
		}
	}

	return repo
}

func setYumGroupRepositoryToResourceData(repo *repository.YumGroupRepository, resourceData *schema.ResourceData) error {
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

func resourceYumGroupRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getYumGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Yum.Group.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceYumGroupRepositoryRead(resourceData, m)
}

func resourceYumGroupRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Group.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setYumGroupRepositoryToResourceData(repo, resourceData)
}

func resourceYumGroupRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getYumGroupRepositoryFromResourceData(resourceData)

	if err := client.Repository.Yum.Group.Update(repoName, repo); err != nil {
		return err
	}

	return resourceYumGroupRepositoryRead(resourceData, m)
}

func resourceYumGroupRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Yum.Group.Delete(resourceData.Id())
}

func resourceYumGroupRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Group.Get(resourceData.Id())
	return repo != nil, err
}
