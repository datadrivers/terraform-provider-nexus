package nexus

import (
	"encoding/json"
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"io/ioutil"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,
		Exists: resourceRepositoryExists,

		Schema: map[string]*schema.Schema{
			"format": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"apt", "bower", "docker", "maven2", "nuget", "pypi"}, false),
			},
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Default:     true,
				Description: "Whether this repository accepts incoming requests",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"group", "hosted", "proxy"}, false),
			},
			"cleanup": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Type:     schema.TypeSet,
						},
					},
				},
			},
			"apt": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"docker"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distribution": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"apt_signing": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"docker"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair": {
							Type:     schema.TypeString,
							Required: true,
						},
						"passphrase": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"bower": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"apt", "apt_signing", "docker"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rewrite_package_urls": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"docker": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"apt", "apt_signing", "bower"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"force_basic_auth": {
							Default:     true,
							Description: "Whether to force authentication (Docker Bearer Token Realm required if false)",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"http_port": {
							Description: "Create an HTTP connector at specified port",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						"https_port": {
							Description: "Create an HTTPS connector at specified port",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						"v1enabled": {
							Default:     false,
							Description: "Whether to allow clients to use the V1 API to interact with this repository",
							Optional:    true,
							Type:        schema.TypeBool,
						},
					},
				},
			},
			"docker_proxy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_type": {
							Description: "Type of Docker Index",
							Required:    true,
							Type:        schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"CUSTOM",
								"HUB",
								"REGISTRY",
							}, false),
						},
						"index_url": {
							Description: "URL of Docker Index to use",
							Required:    true,
							Type:        schema.TypeString,
						},
					},
				},
			},
			"http_client": {
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_block": {
							Default:     true,
							Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"blocked": {
							Default:     false,
							Description: "Whether to block outbound connections on the repository",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"connection": {
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retries": {
										Default:      0,
										Description:  "Total retries if the initial connection attempt suffers a timeout",
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(0, 10),
									},
									"timeout": {
										Default:      60,
										Description:  "Seconds to wait for activity before stopping and retrying the connection",
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(1, 3600),
									},
								},
							},
							Type:     schema.TypeList,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"negative_cache": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Default:     false,
							Description: "Whether to cache responses for content not present in the proxied repository",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"ttl": {
							Default:     1440,
							Description: "How long to cache the fact that a file was not found in the repository (in minutes)",
							Optional:    true,
							Type:        schema.TypeInt,
						},
					},
				},
			},
			"proxy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_max_age": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1440,
						},
						"metadata_max_age": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1440,
						},
						"remote_url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"storage": {
				DefaultFunc: repositoryStorageDefault,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_store_name": {
							Default:     "default",
							Description: "Blob store used to store repository contents",
							Optional:    true,
							Type:        schema.TypeString,
						},
						"strict_content_type_validation": {
							Default:     true,
							Description: "Whether to validate uploaded content's MIME type appropriate for the repository format",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"write_policy": {
							Default:     "ALLOW_ONCE",
							Description: "Controls if deployments of and updates to assets are allowed",
							Optional:    true,
							Type:        schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"ALLOW",
								"ALLOW_ONCE",
								"DENY",
							}, false),
						},
					},
				},
			},
		},
	}
}

func repositoryStorageDefault() (interface{}, error) {
	data := map[string]interface{}{
		"blob_store_name":                "default",
		"strict_content_type_validation": true,
		"write_policy":                   "ALLOW",
	}
	return []map[string]interface{}{data}, nil
}
func getRepositoryFromResourceData(d *schema.ResourceData) nexus.Repository {
	repo := nexus.Repository{
		Name:   d.Get("name").(string),
		Type:   d.Get("type").(string),
		Format: d.Get("format").(string),
		Online: d.Get("online").(bool),
	}

	if _, ok := d.GetOk("apt"); ok {
		aptList := d.Get("apt").([]interface{})
		aptConfig := aptList[0].(map[string]interface{})

		repo.RepositoryApt = &nexus.RepositoryApt{
			Distribution: aptConfig["distribution"].(string),
		}
	}

	if _, ok := d.GetOk("apt_signing"); ok {
		aptSigningList := d.Get("apt_signing").([]interface{})
		aptSigningConfig := aptSigningList[0].(map[string]interface{})

		repo.RepositoryAptSigning = &nexus.RepositoryAptSigning{
			Keypair:    aptSigningConfig["keypair"].(string),
			Passphrase: aptSigningConfig["passphrase"].(string),
		}
	}

	if _, ok := d.GetOk("cleanup"); ok {
		cleanupList := d.Get("cleanup").([]interface{})
		cleanupConfig := cleanupList[0].(map[string]interface{})
		repoCleanupPolicyNames := cleanupConfig["policy_names"].(*schema.Set)

		cleanupPolicyNames := make([]string, repoCleanupPolicyNames.Len())
		for _, v := range repoCleanupPolicyNames.List() {
			cleanupPolicyNames = append(cleanupPolicyNames, v.(string))
		}

		repo.RepositoryCleanup = &nexus.RepositoryCleanup{
			PolicyNames: cleanupPolicyNames,
		}
	}

	if _, ok := d.GetOk("docker"); ok {
		dockerList := d.Get("docker").([]interface{})
		dockerConfig := dockerList[0].(map[string]interface{})
		repo.RepositoryDocker = &nexus.RepositoryDocker{
			ForceBasicAuth: dockerConfig["force_basic_auth"].(bool),
			V1Enabled:      dockerConfig["v1enabled"].(bool),
		}

		if v, ok := dockerConfig["http_port"]; ok {
			port := new(int)
			*port = v.(int)
			repo.RepositoryDocker.HTTPPort = port
		}

		if v, ok := dockerConfig["https_port"]; ok {
			port := new(int)
			*port = v.(int)
			repo.RepositoryDocker.HTTPSPort = port
		}
	}

	if _, ok := d.GetOk("docker_proxy"); ok {
		dockerProxyList := d.Get("docker_proxy").([]interface{})
		dockerProxyConfig := dockerProxyList[0].(map[string]interface{})

		repo.RepositoryDockerProxy = &nexus.RepositoryDockerProxy{
			IndexType: dockerProxyConfig["index_type"].(string),
			IndexURL:  dockerProxyConfig["index_url"].(string),
		}
	}

	if _, ok := d.GetOk("http_client"); ok {
		httpClientList := d.Get("http_client").([]interface{})
		httpClientConfig := httpClientList[0].(map[string]interface{})

		repo.RepositoryHTTPClient = &nexus.RepositoryHTTPClient{
			AutoBlock: httpClientConfig["auto_block"].(bool),
			Blocked:   httpClientConfig["blocked"].(bool),
			// Connection: ,
		}

		if _, ok := httpClientConfig["authentication"]; ok {
			authList := httpClientConfig["authentication"].([]interface{})
			authConfig := authList[0].(map[string]interface{})

			auth := &nexus.RepositoryHTTPClientAuthentication{
				Type:       authConfig["type"].(string),
				Username:   authConfig["username"].(string),
				NTLMDomain: authConfig["ntlm_domain"].(string),
				NTLMHost:   authConfig["ntlm_host"].(string),
			}
			repo.RepositoryHTTPClient.Authentication = *auth
		}
	}

	if _, ok := d.GetOk("negative_cache"); ok {
		negativeCacheList := d.Get("negative_cache").([]interface{})
		negativeCacheConfig := negativeCacheList[0].(map[string]interface{})

		repo.RepositoryNegativeCache = &nexus.RepositoryNegativeCache{
			Enabled: negativeCacheConfig["enabled"].(bool),
			TTL:     negativeCacheConfig["ttl"].(int),
		}
	}

	if _, ok := d.GetOk("proxy"); ok {
		proxyList := d.Get("proxy").([]interface{})
		proxyConfig := proxyList[0].(map[string]interface{})
		repo.RepositoryProxy = &nexus.RepositoryProxy{
			ContentMaxAge:  proxyConfig["content_max_age"].(int),
			MetadataMaxAge: proxyConfig["metadata_max_age"].(int),
			RemoteURL:      proxyConfig["remote_url"].(string),
		}
	}

	if _, ok := d.GetOk("storage"); ok {
		storageList := d.Get("storage").([]interface{})
		storageConfig := storageList[0].(map[string]interface{})

		repo.RepositoryStorage = &nexus.RepositoryStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
			WritePolicy:                 storageConfig["write_policy"].(string),
		}
	}

	return repo
}

func setRepositoryToResourceData(repo *nexus.Repository, d *schema.ResourceData) error {
	d.SetId(repo.Name)
	d.Set("format", repo.Format)
	d.Set("name", repo.Name)
	d.Set("online", repo.Online)
	d.Set("type", repo.Type)

	if repo.RepositoryApt != nil {
		if err := d.Set("apt", flattenRepositoryApt(repo.RepositoryApt)); err != nil {
			return err
		}
	}

	if repo.RepositoryAptSigning != nil {
		if err := d.Set("apt_signing", flattenRepositoryAptSigning(repo.RepositoryAptSigning)); err != nil {
			return err
		}
	}

	if repo.RepositoryDocker != nil {
		if err := d.Set("docker", flattenRepositoryDocker(repo.RepositoryDocker)); err != nil {
			return err
		}
	}

	if repo.RepositoryDockerProxy != nil {
		if err := d.Set("docker_proxy", flattenRepositoryDockerProxy(repo.RepositoryDockerProxy)); err != nil {
			return err
		}
	}

	// if repo.RepositoryNegativeCache != nil {
	// 	if err := d.Set("negative_cache", flattenRepositoryNegativeCache(repo.RepositoryNegativeCache)); err != nil {
	// 		return err
	// 	}
	// }

	// if repo.RepositoryProxy != nil {
	// 	if err := d.Set("proxy", flattenRepositoryProxy(repo.RepositoryProxy)); err != nil {
	// 		return err
	// 	}
	// }

	//	if repo.RepositoryStorage != nil {
	if err := d.Set("storage", flattenRepositoryStorage(repo.RepositoryStorage)); err != nil {
		return err
	}
	//	}

	return nil
}

func flattenRepositoryApt(apt *nexus.RepositoryApt) []map[string]interface{} {
	data := map[string]interface{}{
		"distribution": apt.Distribution,
	}

	return []map[string]interface{}{data}
}

func flattenRepositoryAptSigning(aptSigning *nexus.RepositoryAptSigning) []map[string]interface{} {
	data := map[string]interface{}{
		"keypair":    aptSigning.Keypair,
		"passphrase": aptSigning.Passphrase,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryCleanup(cleanup *nexus.RepositoryCleanup) []map[string]interface{} {
	data := map[string]interface{}{
		// "policy_names":
	}

	return []map[string]interface{}{data}
}

func flattenRepositoryDocker(docker *nexus.RepositoryDocker) []map[string]interface{} {
	data := map[string]interface{}{
		"force_basic_auth": docker.ForceBasicAuth,
		"v1enabled":        docker.V1Enabled,
	}

	if docker.HTTPPort != nil {
		data["http_port"] = *docker.HTTPPort
	}
	if docker.HTTPSPort != nil {
		data["https_port"] = *docker.HTTPSPort
	}

	return []map[string]interface{}{data}
}

func flattenRepositoryDockerProxy(dockerProxy *nexus.RepositoryDockerProxy) []map[string]interface{} {
	data := map[string]interface{}{
		"index_type": dockerProxy.IndexType,
		"index_url":  dockerProxy.IndexURL,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryNegativeCache(negativeCache *nexus.RepositoryNegativeCache) []map[string]interface{} {
	data := map[string]interface{}{
		"enabled": negativeCache.Enabled,
		"ttl":     negativeCache.TTL,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryProxy(proxy *nexus.RepositoryProxy) []map[string]interface{} {
	data := map[string]interface{}{
		"content_max_age":  proxy.ContentMaxAge,
		"metadata_max_age": proxy.MetadataMaxAge,
		"remote_url":       proxy.RemoteURL,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryStorage(storage *nexus.RepositoryStorage) []map[string]interface{} {
	data := map[string]interface{}{
		"blob_store_name":                storage.BlobStoreName,
		"strict_content_type_validation": storage.StrictContentTypeValidation,
		"write_policy":                   storage.WritePolicy,
	}

	return []map[string]interface{}{data}
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	repo := getRepositoryFromResourceData(d)
	repoFormat := d.Get("format").(string)
	repoType := d.Get("type").(string)

	if err := client.RepositoryCreate(repo, repoFormat, repoType); err != nil {
		return err
	}

	if err := setRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	id := d.Id()

	repo, err := nexusClient.RepositoryRead(id)
	if err != nil {
		return err
	}

	data, err := json.Marshal(repo)

	if err := ioutil.WriteFile("/tmp/test.log", data, 0644); err != nil {
		return err
	}

	if repo == nil {
		d.SetId("")
		return nil
	}

	return setRepositoryToResourceData(repo, d)
}

func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceRepositoryRead(d, m)
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	id := d.Get("name").(string)

	return nexusClient.RepositoryDelete(id)
}

func resourceRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	id := d.Get("name").(string)

	repo, err := nexusClient.RepositoryRead(id)
	if err != nil {
		return false, nil
	}

	return repo != nil, nil
}
