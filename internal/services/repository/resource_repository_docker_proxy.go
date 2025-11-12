package repository

import (
	"regexp"
	"strings"

	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	repositorySchema "github.com/datadrivers/terraform-provider-nexus/internal/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceRepositoryDockerProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a docker proxy repository.",

		Create: resourceDockerProxyRepositoryCreate,
		Delete: resourceDockerProxyRepositoryDelete,
		Exists: resourceDockerProxyRepositoryExists,
		Read:   resourceDockerProxyRepositoryRead,
		Update: resourceDockerProxyRepositoryUpdate,
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
			// Docker proxy schemas
			"docker": repositorySchema.ResourceDocker,
			"docker_proxy": {
				Description: "docker_proxy contains the configuration of the docker index",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_type": {
							Description:  "Type of Docker Index. Possible values: `HUB`, `REGISTRY` or `CUSTOM`",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{string(repository.DockerProxyIndexTypeHub), string(repository.DockerProxyIndexTypeRegistry), string(repository.DockerProxyIndexTypeCustom)}, false),
						},
						"index_url": {
							Description:  "Url of Docker Index to use",
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("http[s]?://.*"), "index_url should be in the format 'http://www.example.com'"),
						},
						"cache_foreign_layers": {
							Description: "Allow Nexus Repository Manager to download and cache foreign layers",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"foreign_layer_url_whitelist": {
							Description: "A set of regular expressions used to identify URLs that are allowed for foreign layer requests",
							Optional:    true,
							Type:        schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func getDockerProxyRepositoryFromResourceData(resourceData *schema.ResourceData) repository.DockerProxyRepository {
	httpClientConfig := resourceData.Get("http_client").([]interface{})[0].(map[string]interface{})
	negativeCacheConfig := resourceData.Get("negative_cache").([]interface{})[0].(map[string]interface{})
	proxyConfig := resourceData.Get("proxy").([]interface{})[0].(map[string]interface{})
	storageConfig := resourceData.Get("storage").([]interface{})[0].(map[string]interface{})
	dockerConfig := resourceData.Get("docker").([]interface{})[0].(map[string]interface{})
	dockerProxyConfig := resourceData.Get("docker_proxy").([]interface{})[0].(map[string]interface{})

	repo := repository.DockerProxyRepository{
		Name:   resourceData.Get("name").(string),
		Online: resourceData.Get("online").(bool),
		Storage: repository.Storage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		},
		Docker: repository.Docker{
			ForceBasicAuth: dockerConfig["force_basic_auth"].(bool),
			V1Enabled:      dockerConfig["v1_enabled"].(bool),
		},
		DockerProxy: repository.DockerProxy{
			IndexType: repository.DockerProxyIndexType(dockerProxyConfig["index_type"].(string)),
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

	if httpPort, ok := dockerConfig["http_port"]; ok {
		if httpPort.(int) > 0 {
			repo.Docker.HTTPPort = tools.GetIntPointer(httpPort.(int))
		}
	}

	if httpsPort, ok := dockerConfig["https_port"]; ok {
		if httpsPort.(int) > 0 {
			repo.Docker.HTTPSPort = tools.GetIntPointer(httpsPort.(int))
		}
	}

	if subdomain, ok := dockerConfig["subdomain"]; ok {
		if subdomain.(string) != "" {
			repo.Docker.Subdomain = tools.GetStringPointer(subdomain.(string))
		}
	}

	if dockerProxyConfig["index_url"].(string) != "" {
		repo.DockerProxy.IndexURL = tools.GetStringPointer(strings.TrimSpace(dockerProxyConfig["index_url"].(string)))
	}

	cacheForeignLayers, ok := dockerProxyConfig["cache_foreign_layers"]
	if ok {
		repo.DockerProxy.CacheForeignLayers = tools.GetBoolPointer(cacheForeignLayers.(bool))
	}

	foreignLayerUrlWhitelist, ok := dockerProxyConfig["foreign_layer_url_whitelist"]
	if ok {
		repo.DockerProxy.ForeignLayerUrlWhitelist = tools.InterfaceSliceToStringSlice(foreignLayerUrlWhitelist.(*schema.Set).List())
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

func setDockerProxyRepositoryToResourceData(repo *repository.DockerProxyRepository, resourceData *schema.ResourceData) error {
	resourceData.SetId(repo.Name)
	resourceData.Set("name", repo.Name)
	resourceData.Set("online", repo.Online)

	if err := resourceData.Set("docker", flattenDocker(&repo.Docker)); err != nil {
		return err
	}

	if err := resourceData.Set("docker_proxy", flattenDockerProxy(&repo.DockerProxy)); err != nil {
		return err
	}

	if repo.RoutingRuleName != nil {
		resourceData.Set("routing_rule", repo.RoutingRuleName)
	} else if repo.RoutingRule != nil {
		resourceData.Set("routing_rule", repo.RoutingRule)
	} else if repo.RoutingRuleName == nil && repo.RoutingRule == nil {
		resourceData.Set("routing_rule", nil)
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

func resourceDockerProxyRepositoryCreate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getDockerProxyRepositoryFromResourceData(resourceData)

	if err := client.Repository.Docker.Proxy.Create(repo); err != nil {
		return err
	}
	resourceData.SetId(repo.Name)

	return resourceDockerProxyRepositoryRead(resourceData, m)
}

func resourceDockerProxyRepositoryRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Docker.Proxy.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		resourceData.SetId("")
		return nil
	}

	return setDockerProxyRepositoryToResourceData(repo, resourceData)
}

func resourceDockerProxyRepositoryUpdate(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := resourceData.Id()
	repo := getDockerProxyRepositoryFromResourceData(resourceData)

	if err := client.Repository.Docker.Proxy.Update(repoName, repo); err != nil {
		return err
	}

	return resourceDockerProxyRepositoryRead(resourceData, m)
}

func resourceDockerProxyRepositoryDelete(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	return client.Repository.Docker.Proxy.Delete(resourceData.Id())
}

func resourceDockerProxyRepositoryExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Docker.Proxy.Get(resourceData.Id())
	return repo != nil, err
}
