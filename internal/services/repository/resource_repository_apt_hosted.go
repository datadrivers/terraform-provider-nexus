package repository

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/dre2004/terraform-provider-nexus/internal/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryAptHosted() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a hosted apt repository.",

		Create: resourceAptHostedRepositoryCreate,
		Delete: resourceAptHostedRepositoryDelete,
		Exists: resourceAptHostedRepositoryExists,
		Read:   resourceAptHostedRepositoryRead,
		Update: resourceAptHostedRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.ResourceID,
			"name":   repositorySchema.ResourceName,
			"online": repositorySchema.ResourceOnline,
			// Hosted schemas
			"cleanup":   repositorySchema.ResourceCleanup,
			"component": repositorySchema.ResourceComponent,
			"storage":   repositorySchema.ResourceHostedStorage,
			// Apt hosted schemas
			"distribution": {
				Description: "Distribution to fetch",
				Required:    true,
				Type:        schema.TypeString,
			},
			"signing": {
				Description: "Signing contains signing data of hosted repositores of format Apt",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair": {
							Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
						"passphrase": {
							Description: "Passphrase to access PGP signing key",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func getAptHostedRepositoryFromResourceData(resourceData *schema.ResourceData) repository.AptHostedRepository {
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))
	signingConfig := resourceData.Get("signing").([]interface{})[0].(map[string]interface{})

	repo := repository.AptHostedRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.HostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 &writePolicy,
		},
		Apt: repository.AptHosted{
			Distribution: resourceData.Get("distribution").(string),
		},
		AptSigning: repository.AptSigning{
			Keypair: signingConfig["keypair"].(string),
		},
	}

	if signingConfig["passphrase"] != nil {
		repo.AptSigning.Passphrase = tools.GetStringPointer(signingConfig["passphrase"].(string))
	}

	cleanupList := resourceData.Get("cleanup").([]interface{})
	if len(cleanupList) > 0 && cleanupList[0] != nil {
		cleanupConfig := cleanupList[0].(map[string]interface{})
		if len(cleanupConfig) > 0 {
			policy_names, ok := cleanupConfig["policy_names"]
			if ok {
				repo.Cleanup = &repository.Cleanup{
					PolicyNames: tools.InterfaceSliceToStringSlice(policy_names.(*schema.Set).List()),
				}
			}
		}
	}

	componentList := resourceData.Get("component").([]interface{})
	if len(componentList) > 0 && componentList[0] != nil {
		componentConfig := componentList[0].(map[string]interface{})
		if len(componentConfig) > 0 {
			repo.Component = &repository.Component{
				ProprietaryComponents: componentConfig["proprietary_components"].(bool),
			}
		}
	}

	return repo
}

func setAptHostedRepositoryToResourceData(repo *repository.AptHostedRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)
	resourceData.Set("distribution", repo.Apt.Distribution)

	if err := resourceData.Set("storage", flattenHostedStorage(&repo.Storage)); err != nil {
		return err
	}

	if repo.Cleanup != nil {
		if err := resourceData.Set("cleanup", flattenCleanup(repo.Cleanup)); err != nil {
			return err
		}
	}

	if repo.Component != nil {
		if err := resourceData.Set("component", flattenComponent(repo.Component)); err != nil {
			return err
		}
	}

	return nil
}

func resourceAptHostedRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getAptHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Apt.Hosted.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceAptHostedRepositoryRead(resourceData, m)
}

func resourceAptHostedRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Apt.Hosted.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setAptHostedRepositoryToResourceData(repo, resourceData)
}

func resourceAptHostedRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getAptHostedRepositoryFromResourceData(resourceData)

	if err := client.Repository.Apt.Hosted.Update(repoName, repo); err != nil {
		return err
	}

	return resourceAptHostedRepositoryRead(resourceData, m)
}

func resourceAptHostedRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Apt.Hosted.Delete(resourceData.Id())
}

func resourceAptHostedRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Apt.Hosted.Get(resourceData.Id())
	return repo != nil, err
}
