/*
Use this resource to create a Nexus Repository.

Example Usage

Example Usage - Apt hosted repository

```hcl
resource "nexus_repository" "apt_hosted" {
  name   = "apt-repo"
  format = "apt"
  type   = "hosted"

  apt {
    distribution = "bionic"
  }

  apt_signing {
    keypair    = "<keypair>"
    passphrase = "<passphrase>"
  }

  storage {
    blob_store_name                = "default"
    strict_content_type_validation = true
    write_policy                   = "ALLOW_ONCE"
  }
}
```

Example Usage - Docker group repository

```hcl
resource "nexus_repository" "docker_group" {
  name   = "docker-group"
  format = "docker"
  type   = "group"
  online = true

  group {
    member_names = [
      "docker_releases",
      "docker_hub"
    ]
  }

  docker {
    force_basic_auth = false
    http_port        = 5000
    https_port       = 0
    v1enabled        = false
  }

  storage {
    blob_store_name                = "docker_group_blobstore"
    strict_content_type_validation = true
  }
}
```
*/
package nexus

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,
		Exists: resourceRepositoryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"format": {
				Description:  "Repository format. Possible values: `apt`, `bower`, `conan`, `docker`, `gitlfs`, `go`, `helm`, `maven2`, `npm`, `nuget`, `p2`, `pypi`, `raw`, `rubygems`, `yum`",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(repository.RepositoryFormats, false),
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
				Description:  "Repository type. Possible values: `group`, `hosted`, `proxy`",
				ForceNew:     true,
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(repository.RepositoryTypes, false),
			},
			"apt": {
				Description:   "Apt specific configuration of the repository",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"bower", "docker", "docker_proxy", "maven"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distribution": {
							Description: "The linux distribution",
							Type:        schema.TypeString,
							Required:    true,
						},
						"flat": {
							Description: "Whether this repository is flat",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			"apt_signing": {
				Description:   "Apt signing configuration for the repository",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"bower", "docker", "docker_proxy", "maven"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keypair": {
							Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor )",
							Type:        schema.TypeString,
							Required:    true,
						},
						"passphrase": {
							Description: "Passphrase for the keypair",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"bower": {
				Description:   "Bower specific configuration of the repository",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"apt", "apt_signing", "docker", "docker_proxy", "maven"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rewrite_package_urls": {
							Description: "Force Bower to retrieve packages through the proxy repository",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"cleanup": {
				Description: "Cleanup policies",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Description: "List of policy names",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Set: func(v interface{}) int {
								return schema.HashString(strings.ToLower(v.(string)))
							},
							Type: schema.TypeSet,
						},
					},
				},
			},
			"docker": {
				Description:   "Docker specific configuration of the repository",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"apt", "apt_signing", "bower", "maven"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"force_basic_auth": {
							Default:     true,
							Description: "Whether to force authentication (Docker Bearer Token Realm required if false)",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"http_port": {
							Default:     0,
							Description: "Create an HTTP connector at specified port",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						"https_port": {
							Default:     0,
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
				ConflictsWith: []string{"apt", "apt_signing", "bower", "maven"},
				Description:   "Configuration for docker proxy repository",
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
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
							Optional:    true,
							Type:        schema.TypeString,
						},
					},
				},
			},
			"group": {
				Description: "Configuration for repository group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_names": {
							Description: "Member repositories names",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
							Set: func(v interface{}) int {
								return schema.HashString(strings.ToLower(v.(string)))
							},
							Type: schema.TypeSet,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"http_client": {
				Description: "HTTP Client configuration for proxy repositories. Required for docker proxy repositories.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Description: "Authentication configuration of the HTTP client",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description:  "Authentication type. Possible values: `ntlm`, `username` or `bearerToken`. Only npm supports bearerToken authentication",
										Optional:     true,
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"ntlm", "username", "bearerToken"}, false),
									},
									"username": {
										Description: "The username used by the proxy repository",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"password": {
										Description: "The password used by the proxy repository",
										Optional:    true,
										Sensitive:   true,
										Type:        schema.TypeString,
									},
									"ntlm_domain": {
										Description: "The ntlm domain to connect",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"ntlm_host": {
										Description: "The ntlm host to connect",
										Optional:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
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
							Description: "Connection configuration of the HTTP client",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_cookies": {
										Description: "Whether to allow cookies to be stored and used",
										Optional:    true,
										Type:        schema.TypeBool,
									},
									"retries": {
										Description:  "Total retries if the initial connection attempt suffers a timeout",
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(0, 10),
									},
									"timeout": {
										Description:  "Seconds to wait for activity before stopping and retrying the connection",
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(1, 3600),
									},
									"user_agent_suffix": {
										Description: "Custom fragment to append to User-Agent header in HTTP requests",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"use_trust_store": {
										Description: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
										Optional:    true,
										Default:     false,
										Type:        schema.TypeBool,
									},
								},
							},
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"maven": {
				Description: "Maven specific configuration of the repository",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_policy": {
							Description: "What type of artifacts does this repository store? Possible values: `RELEASE`, `SNAPSHOT` or `MIXED`",
							Default:     "RELEASE",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"layout_policy": {
							Description: "Validate that all paths are maven artifact or metadata paths. Possible values: `PERMISSIVE` or `STRICT`",
							Default:     "PERMISSIVE",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"negative_cache": {
				Description: "Configuration of the negative cache handling",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
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
			"nuget_proxy": {
				Description: "Configuration for the nuget proxy repository",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_cache_item_max_age": {
							Description: "What type of artifacts does this repository store",
							Required:    true,
							Type:        schema.TypeInt,
						},
						"nuget_version": {
							Description:  "Nuget protocol version. Possible values: `V2` or `V3` (Default)",
							Default:      "V3",
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"V2", "V3"}, false),
							Optional:     true,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"proxy": {
				Description: "Configuration for the proxy repository",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_max_age": {
							Description: "How long (in minutes) to cache artifacts before rechecking the remote repository",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1440,
						},
						"metadata_max_age": {
							Description: "How long (in minutes) to cache metadata before rechecking the remote repository.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1440,
						},
						"remote_url": {
							Description: "Location of the remote repository being proxied",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"storage": {
				Description: "The storage configuration of the repository",
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
							Description: "Controls if deployments of and updates to assets are allowed. Possible values: `ALLOW`, `ALLOW_ONCE`, `DENY`",
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
			"yum": {
				Description:   "Yum specific configuration of the repository",
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"bower", "docker", "docker_proxy", "maven", "apt"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repodata_depth": {
							Description: "Specifies the repository depth where repodata folder(s) are created. Possible values: 0-5",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"deploy_policy": {
							Description:  "Validate that all paths are RPMs or yum metadata. Possible values: `STRICT` or `PERMISSIVE`",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"STRICT", "PERMISSIVE"}, false),
						},
					},
				},
			},
		},
	}
}

func RepositoryCleanupDefault() (interface{}, error) {
	data := map[string]interface{}{
		"policy_names": []string{},
	}
	return []map[string]interface{}{data}, nil
}

func repositoryStorageDefault() (interface{}, error) {

	data := map[string]interface{}{
		"blob_store_name":                "default",
		"strict_content_type_validation": true,
	}
	return []map[string]interface{}{data}, nil
}

func getRepositoryFromResourceData(d *schema.ResourceData) repository.LegacyRepository {
	repo := repository.LegacyRepository{
		Format: d.Get("format").(string),
		Name:   d.Get("name").(string),
		Online: d.Get("online").(bool),
		Type:   d.Get("type").(string),
	}

	if _, ok := d.GetOk("apt"); ok {
		aptList := d.Get("apt").([]interface{})
		aptConfig := aptList[0].(map[string]interface{})

		repo.Apt = &repository.AptProxy{
			Distribution: aptConfig["distribution"].(string),
			Flat:         aptConfig["flat"].(bool),
		}
	}

	if _, ok := d.GetOk("apt_signing"); ok {
		aptSigningList := d.Get("apt_signing").([]interface{})
		aptSigningConfig := aptSigningList[0].(map[string]interface{})
		keypair := aptSigningConfig["keypair"].(string)
		passphrase := aptSigningConfig["passphrase"].(string)

		repo.AptSigning = &repository.AptSigning{
			Keypair:    &keypair,
			Passphrase: &passphrase,
		}
	}

	if _, ok := d.GetOk("bower"); ok {
		bowerList := d.Get("bower").([]interface{})
		bowerConfig := bowerList[0].(map[string]interface{})

		repo.Bower = &repository.Bower{
			RewritePackageUrls: bowerConfig["rewrite_package_urls"].(bool),
		}
	}

	if _, ok := d.GetOk("cleanup"); ok {
		cleanupList := d.Get("cleanup").([]interface{})
		cleanupConfig := cleanupList[0].(map[string]interface{})
		repo.Cleanup = &repository.Cleanup{
			PolicyNames: interfaceSliceToStringSlice(cleanupConfig["policy_names"].(*schema.Set).List()),
		}
	}

	if _, ok := d.GetOk("docker"); ok {
		dockerList := d.Get("docker").([]interface{})
		dockerConfig := dockerList[0].(map[string]interface{})
		repo.Docker = &repository.Docker{
			ForceBasicAuth: dockerConfig["force_basic_auth"].(bool),
			V1Enabled:      dockerConfig["v1enabled"].(bool),
		}

		if v, ok := dockerConfig["http_port"]; ok {
			value := v.(int)
			if value > 0 {
				port := new(int)
				*port = value
				repo.Docker.HTTPPort = port
			}
		}

		if v, ok := dockerConfig["https_port"]; ok {
			value := v.(int)
			if value > 0 {
				port := new(int)
				*port = v.(int)
				repo.Docker.HTTPSPort = port
			}
		}
	}

	if _, ok := d.GetOk("docker_proxy"); ok {
		dockerProxyList := d.Get("docker_proxy").([]interface{})
		dockerProxyConfig := dockerProxyList[0].(map[string]interface{})

		var indexURLValue *string
		indexURL := strings.TrimSpace(dockerProxyConfig["index_url"].(string))
		if indexURL != "" {
			indexURLValue = &indexURL
		}
		indexType := repository.DockerProxyIndexType(dockerProxyConfig["index_type"].(string))
		repo.DockerProxy = &repository.DockerProxy{
			IndexType: &indexType,
			IndexURL:  indexURLValue,
		}
	}

	if _, ok := d.GetOk("group"); ok {
		groupList := d.Get("group").([]interface{})
		groupMemberNames := make([]string, 0)

		if len(groupList) == 1 && groupList[0] != nil {
			groupConfig := groupList[0].(map[string]interface{})
			groupConfigMemberNames := groupConfig["member_names"].(*schema.Set)

			for _, v := range groupConfigMemberNames.List() {
				groupMemberNames = append(groupMemberNames, v.(string))
			}
		}
		repo.Group = &repository.Group{
			MemberNames: groupMemberNames,
		}
	}

	if _, ok := d.GetOk("http_client"); ok {
		httpClientList := d.Get("http_client").([]interface{})
		httpClientConfig := httpClientList[0].(map[string]interface{})

		repo.HTTPClient = &repository.HTTPClient{
			AutoBlock: httpClientConfig["auto_block"].(bool),
			Blocked:   httpClientConfig["blocked"].(bool),
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
			connList := v.([]interface{})
			if len(connList) == 1 && connList[0] != nil {
				connConfig := connList[0].(map[string]interface{})

				repo.HTTPClient.Connection = &repository.HTTPClientConnection{
					EnableCookies:   connConfig["enable_cookies"].(*bool),
					Retries:         connConfig["retries"].(*int),
					Timeout:         connConfig["timeout"].(*int),
					UserAgentSuffix: connConfig["user_agent_suffix"].(string),
					UseTrustStore:   connConfig["use_trust_store"].(*bool),
				}
			}
		}
	}

	if _, ok := d.GetOk("maven"); ok {
		mavenList := d.Get("maven").([]interface{})
		mavenConfig := mavenList[0].(map[string]interface{})

		versionPolicy := repository.MavenVersionPolicy(mavenConfig["version_policy"].(string))
		layoutPolicy := repository.MavenLayoutPolicy(mavenConfig["layout_policy"].(string))
		repo.Maven = &repository.Maven{
			VersionPolicy: &versionPolicy,
			LayoutPolicy:  &layoutPolicy,
		}
	}

	if _, ok := d.GetOk("negative_cache"); ok {
		negativeCacheList := d.Get("negative_cache").([]interface{})
		negativeCacheConfig := negativeCacheList[0].(map[string]interface{})

		repo.NegativeCache = &repository.NegativeCache{
			Enabled: negativeCacheConfig["enabled"].(bool),
			TTL:     negativeCacheConfig["ttl"].(int),
		}
	}

	if _, ok := d.GetOk("nuget_proxy"); ok {
		nugetProxyList := d.Get("nuget_proxy").([]interface{})
		nugetProxyConfig := nugetProxyList[0].(map[string]interface{})

		repo.NugetProxy = &repository.NugetProxy{
			QueryCacheItemMaxAge: nugetProxyConfig["query_cache_item_max_age"].(int),
			NugetVersion:         repository.NugetVersion(nugetProxyConfig["nuget_version"].(string)),
		}
	}

	if _, ok := d.GetOk("proxy"); ok {
		proxyList := d.Get("proxy").([]interface{})
		proxyConfig := proxyList[0].(map[string]interface{})
		remoteURL := proxyConfig["remote_url"].(string)
		repo.Proxy = &repository.Proxy{
			ContentMaxAge:  proxyConfig["content_max_age"].(int),
			MetadataMaxAge: proxyConfig["metadata_max_age"].(int),
			RemoteURL:      &remoteURL,
		}
	}

	if _, ok := d.GetOk("storage"); ok {
		storageList := d.Get("storage").([]interface{})
		storageConfig := storageList[0].(map[string]interface{})

		repo.Storage = &repository.HostedStorage{
			BlobStoreName:               storageConfig["blob_store_name"].(string),
			StrictContentTypeValidation: storageConfig["strict_content_type_validation"].(bool),
		}
		// Only hosted repository has attribute WritePolicy
		if repo.Type == repository.RepositoryTypeHosted {
			writePolicy := repository.StorageWritePolicy(storageConfig["write_policy"].(string))
			repo.Storage.WritePolicy = &writePolicy
		}
	}

	if _, ok := d.GetOk("yum"); ok {
		yumList := d.Get("yum").([]interface{})
		yumConfig := yumList[0].(map[string]interface{})
		deployPolicy := repository.YumDeployPolicy(yumConfig["deploy_policy"].(string))
		repo.Yum = &repository.Yum{
			RepodataDepth: yumConfig["repodata_depth"].(int),
			DeployPolicy:  &deployPolicy,
		}
	}

	return repo
}

func setRepositoryToResourceData(repo *repository.LegacyRepository, d *schema.ResourceData) error {
	d.SetId(repo.Name)
	d.Set("format", repo.Format)
	d.Set("name", repo.Name)
	d.Set("online", repo.Online)
	d.Set("type", repo.Type)

	if repo.Apt != nil {
		if err := d.Set("apt", flattenRepositoryApt(repo.Apt)); err != nil {
			return err
		}
	}

	if repo.AptSigning != nil {
		if err := d.Set("apt_signing", flattenRepositoryAptSigning(repo.AptSigning)); err != nil {
			return err
		}
	}

	if repo.Bower != nil {
		if err := d.Set("bower", flattenRepositoryBower(repo.Bower)); err != nil {
			return err
		}
	}

	if repo.Cleanup != nil {
		if err := d.Set("cleanup", flattenRepositoryCleanup(repo.Cleanup)); err != nil {
			return err
		}
	}

	if repo.Docker != nil {
		if err := d.Set("docker", flattenRepositoryDocker(repo.Docker)); err != nil {
			return err
		}
	}

	if repo.DockerProxy != nil {
		if err := d.Set("docker_proxy", flattenRepositoryDockerProxy(repo.DockerProxy)); err != nil {
			return err
		}
	}

	if repo.Group != nil {
		if err := d.Set("group", flattenRepositoryGroup(repo.Group)); err != nil {
			return err
		}
	}

	if repo.HTTPClient != nil {
		if err := d.Set("http_client", flattenRepositoryHTTPClient(repo.HTTPClient, d)); err != nil {
			return err
		}
	}

	if repo.Maven != nil {
		if err := d.Set("maven", flattenRepositoryMaven(repo.Maven)); err != nil {
			return err
		}
	}

	if repo.NegativeCache != nil {
		if err := d.Set("negative_cache", flattenRepositoryNegativeCache(repo.NegativeCache)); err != nil {
			return err
		}
	}

	if repo.NugetProxy != nil {
		if err := d.Set("nuget_proxy", flattenRepositoryNugetProxy(repo.NugetProxy)); err != nil {
			return err
		}
	}

	if repo.Proxy != nil {
		if err := d.Set("proxy", flattenRepositoryProxy(repo.Proxy)); err != nil {
			return err
		}
	}

	if repo.Yum != nil {
		if err := d.Set("yum", flattenRepositoryYum(repo.Yum)); err != nil {
			return err
		}
	}

	if err := d.Set("storage", flattenRepositoryStorage(repo.Storage, d)); err != nil {
		return err
	}

	return nil
}

func flattenRepositoryApt(apt *repository.AptProxy) []map[string]interface{} {
	if apt == nil {
		return nil
	}
	data := map[string]interface{}{
		"distribution": apt.Distribution,
		"flat":         apt.Flat,
	}

	return []map[string]interface{}{data}
}

func flattenRepositoryAptSigning(aptSigning *repository.AptSigning) []map[string]interface{} {
	if aptSigning == nil {
		return nil
	}
	data := map[string]interface{}{
		"keypair":    aptSigning.Keypair,
		"passphrase": aptSigning.Passphrase,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryBower(bower *repository.Bower) []map[string]interface{} {
	if bower == nil {
		return nil
	}
	data := map[string]interface{}{
		"rewrite_package_urls": bower.RewritePackageUrls,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryCleanup(cleanup *repository.Cleanup) []map[string]interface{} {
	if cleanup == nil {
		return nil
	}
	data := map[string]interface{}{
		"policy_names": stringSliceToInterfaceSlice(cleanup.PolicyNames),
	}

	return []map[string]interface{}{data}
}

func flattenRepositoryDocker(docker *repository.Docker) []map[string]interface{} {
	if docker == nil {
		return nil
	}
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

func flattenRepositoryDockerProxy(dockerProxy *repository.DockerProxy) []map[string]interface{} {
	if dockerProxy == nil {
		return nil
	}
	data := map[string]interface{}{
		"index_type": dockerProxy.IndexType,
		"index_url":  dockerProxy.IndexURL,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryGroup(group *repository.Group) []map[string]interface{} {
	if group == nil {
		return nil
	}
	data := map[string]interface{}{
		"member_names": stringSliceToInterfaceSlice(group.MemberNames),
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryHTTPClient(httpClient *repository.HTTPClient, d *schema.ResourceData) []map[string]interface{} {
	if httpClient == nil {
		return nil
	}
	data := map[string]interface{}{
		"authentication": flattenRepositoryHTTPClientAuthentication(httpClient.Authentication, d),
		"auto_block":     httpClient.AutoBlock,
		"blocked":        httpClient.Blocked,
		// "connection":     flattenRepositoryHTTPClientConnection(httpClient.Connection),
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryHTTPClientAuthentication(auth *repository.HTTPClientAuthentication, d *schema.ResourceData) []map[string]interface{} {
	if auth == nil {
		return nil
	}
	data := map[string]interface{}{
		"ntlm_domain": auth.NTLMDomain,
		"ntlm_host":   auth.NTLMHost,
		"type":        auth.Type,
		"username":    auth.Username,
		"password":    d.Get("http_client.0.authentication.0.password").(string),
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryHTTPClientConnection(conn *repository.HTTPClientConnection) []map[string]interface{} {
	if conn == nil {
		return nil
	}
	data := map[string]interface{}{
		"user_agent_suffix": conn.UserAgentSuffix,
	}
	if conn.EnableCookies != nil {
		data["enable_cookies"] = *conn.EnableCookies
	}
	if conn.Retries != nil {
		data["retries"] = *conn.Retries
	}
	if conn.Timeout != nil {
		data["timeout"] = *conn.Timeout
	}
	if conn.UseTrustStore != nil {
		data["use_trust_store"] = *conn.UseTrustStore
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryMaven(maven *repository.Maven) []map[string]interface{} {
	if maven == nil {
		return nil
	}
	data := map[string]interface{}{
		"version_policy": maven.VersionPolicy,
		"layout_policy":  maven.LayoutPolicy,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryNegativeCache(negativeCache *repository.NegativeCache) []map[string]interface{} {
	if negativeCache == nil {
		return nil
	}
	data := map[string]interface{}{
		"enabled": negativeCache.Enabled,
		"ttl":     negativeCache.TTL,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryNugetProxy(nugetProxy *repository.NugetProxy) []map[string]interface{} {
	if nugetProxy == nil {
		return nil
	}
	data := map[string]interface{}{
		"query_cache_item_max_age": nugetProxy.QueryCacheItemMaxAge,
		"nuget_version":            nugetProxy.NugetVersion,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryProxy(proxy *repository.Proxy) []map[string]interface{} {
	if proxy == nil {
		return nil
	}
	data := map[string]interface{}{
		"content_max_age":  proxy.ContentMaxAge,
		"metadata_max_age": proxy.MetadataMaxAge,
		"remote_url":       proxy.RemoteURL,
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryStorage(storage *repository.HostedStorage, d *schema.ResourceData) []map[string]interface{} {
	if storage == nil {
		return nil
	}
	data := map[string]interface{}{
		"blob_store_name":                storage.BlobStoreName,
		"strict_content_type_validation": storage.StrictContentTypeValidation,
	}
	if d.Get("type") == repository.RepositoryTypeHosted {
		data["write_policy"] = storage.WritePolicy
	}
	return []map[string]interface{}{data}
}

func flattenRepositoryYum(yum *repository.Yum) []map[string]interface{} {
	if yum == nil {
		return nil
	}
	data := map[string]interface{}{
		"deploy_policy":  yum.DeployPolicy,
		"repodata_depth": yum.RepodataDepth,
	}
	return []map[string]interface{}{data}
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo := getRepositoryFromResourceData(d)

	if err := client.Repository.Legacy.Create(repo); err != nil {
		return err
	}

	if err := setRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Legacy.Get(d.Id())
	if err != nil {
		return err
	}

	if repo == nil {
		d.SetId("")
		return nil
	}

	return setRepositoryToResourceData(repo, d)
}

func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	repoName := d.Id()
	repo := getRepositoryFromResourceData(d)

	if err := client.Repository.Legacy.Update(repoName, repo); err != nil {
		return err
	}

	if err := setRepositoryToResourceData(&repo, d); err != nil {
		return err
	}

	return resourceRepositoryRead(d, m)
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	return client.Repository.Legacy.Delete(d.Id())
}

func resourceRepositoryExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	repo, err := client.Repository.Legacy.Get(d.Id())
	return repo != nil, err
}
