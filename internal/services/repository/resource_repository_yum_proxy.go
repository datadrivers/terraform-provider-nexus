package repository

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRepositoryYumProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a yum proxy repository.",

		Create: resourceYumProxyRepositoryCreate,
		Delete: resourceYumProxyRepositoryDelete,
		Exists: resourceYumProxyRepositoryExists,
		Read:   resourceYumProxyRepositoryRead,
		Update: resourceYumProxyRepositoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.ResourceID,
			"name":   repositorySchema.ResourceName,
			"online": repositorySchema.ResourceOnline,
			// Proxy schemas
			"cleanup":        repositorySchema.ResourceCleanup,
			"http_client":    repositorySchema.ResourceHTTPClient,
			"negative_cache": repositorySchema.ResourceNegativeCache,
			"proxy":          repositorySchema.ResourceProxy,
			"routing_rule":   repositorySchema.ResourceRoutingRule,
			"storage":        repositorySchema.ResourceStorage,
			// Yum proxy schemas
			"yum_signing": repositorySchema.ResourceYumSigning,
		},
	}
}

func getYumProxyRepositoryFromResourceData(resourceData *schema.ResourceData) repository.YumProxyRepository {
	httpClientConfig := resourceData.Get("http_client").([]interface{})[0].(map[string]interface{})
	negativeCacheConfig := resourceData.Get("negative_cache").([]interface{})[0].(map[string]interface{})
	proxyConfig := resourceData.Get("proxy").([]interface{})[0].(map[string]interface{})
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})

	repo := repository.YumProxyRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.Storage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		},
		HTTPClient: repository.HTTPClient{
			AutoBlock: httpClientConfig["auto_block"].(bool),
			Blocked:   httpClientConfig["blocked"].(bool),
		},
		NegativeCache: repository.NegativeCache{
			Enabled: negativeCacheConfig["enabled"].(bool),
			TTL:     negativeCacheConfig["ttl"].(int),
		},
		Proxy: repository.Proxy{
			ContentMaxAge:  proxyConfig["content_max_age"].(int),
			MetadataMaxAge: proxyConfig["metadata_max_age"].(int),
			RemoteURL:      proxyConfig["remote_url"].(string),
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

	if routingRule, ok := resourceData.GetOk("routing_rule"); ok {
		repo.RoutingRule = tools.GetStringPointer(routingRule.(string))
		repo.RoutingRuleName = tools.GetStringPointer(routingRule.(string))
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

	if v, ok := httpClientConfig["authentication"]; ok {
		authList := v.([]interface{})
		if len(authList) == 1 && authList[0] != nil {
			authConfig := authList[0].(map[string]interface{})

			repo.HTTPClient.Authentication = &repository.HTTPClientAuthentication{
				NTLMDomain: authConfig["ntlm_domain"].(string),
				NTLMHost:   authConfig["ntlm_host"].(string),
				Type:       repository.HTTPClientAuthenticationType(authConfig["type"].(string)),
				Username:   authConfig["username"].(string),
				Password:   authConfig["password"].(string),
			}
		}
	}

	if v, ok := httpClientConfig["connection"]; ok {
		repo.HTTPClient.Connection = getHTTPClientConnection(v.([]interface{}))
	}

	return repo
}

func setYumProxyRepositoryToResourceData(repo *repository.YumProxyRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if repo.RoutingRuleName != nil {
		resourceData.Set("routing_rule", repo.RoutingRuleName)
	} else if repo.RoutingRule != nil {
		resourceData.Set("routing_rule", repo.RoutingRule)
	}

	if err := resourceData.Set("storage", flattenStorage(&repo.Storage)); err != nil {
		return err
	}

	if err := resourceData.Set("http_client", flattenHTTPClient(&repo.HTTPClient, resourceData)); err != nil {
		return err
	}

	if err := resourceData.Set("negative_cache", flattenNegativeCache(&repo.NegativeCache)); err != nil {
		return err
	}

	if err := resourceData.Set("proxy", flattenProxy(&repo.Proxy)); err != nil {
		return err
	}

	if repo.Cleanup != nil {
		if err := resourceData.Set("cleanup", flattenCleanup(repo.Cleanup)); err != nil {
			return err
		}
	}
	return nil
}

func resourceYumProxyRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getYumProxyRepositoryFromResourceData(resourceData)

	if err := client.Repository.Yum.Proxy.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceYumProxyRepositoryRead(resourceData, m)
}

func resourceYumProxyRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Proxy.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setYumProxyRepositoryToResourceData(repo, resourceData)
}

func resourceYumProxyRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getYumProxyRepositoryFromResourceData(resourceData)

	if err := client.Repository.Yum.Proxy.Update(repoName, repo); err != nil {
		return err
	}

	expectedRoutingRule := resourceData.Get("routing_rule").(string)
	if repo.RoutingRule != nil && *repo.RoutingRule != expectedRoutingRule {
		resourceData.Set("routing_rule", *repo.RoutingRule)
	} else if repo.RoutingRule == nil {
		resourceData.Set("routing_rule", nil)
	}

	return resourceYumProxyRepositoryRead(resourceData, m)
}

func resourceYumProxyRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Yum.Proxy.Delete(resourceData.Id())
}

func resourceYumProxyRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Yum.Proxy.Get(resourceData.Id())
	return repo != nil, err
}
