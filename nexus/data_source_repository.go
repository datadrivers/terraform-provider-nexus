/*
Use this data source to get a repository data structure

Example Usage

```hcl
data "nexus_repository" "maven-central" {
  name = "maven-central"
}
```
*/
package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRepositoryRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"format": {
				Description: "Repository format",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Optional:    true,
				Description: "Whether this repository accepts incoming requests",
				Type:        schema.TypeBool,
			},
			"type": {
				Description: "Repository type",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cleanup": {
				Description: "Cleanup policies",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Description: "List of policy names",
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
						},
					},
				},
			},
			"apt": {
				Description: "Apt specific configuration of the repository",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distribution": {
							Description: "The linux distribution",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"apt_signing": {
				Description: "Apt signing configuration for the repository",
				Type:        schema.TypeList,
				Optional:    true,
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
			"docker": {
				Description: "Docker specific configuration of the repository",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"force_basic_auth": {
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
							Description: "Whether to allow clients to use the V1 API to interact with this repository",
							Optional:    true,
							Type:        schema.TypeBool,
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
							Type:     schema.TypeSet,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"http_client": {
				Description: "HTTP Client configuration for proxy repositories",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Description: "Authentication configuration of the HTTP client",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Authentication type",
										Required:    true,
										Type:        schema.TypeString,
									},
									"username": {
										Description: "The username used by the proxy repository",
										Optional:    true,
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
							Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"blocked": {
							Description: "Whether to block outbound connections on the repository",
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"connection": {
							Description: "Connection configuration of the HTTP client",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retries": {
										Description: "Total retries if the initial connection attempt suffers a timeout",
										Optional:    true,
										Type:        schema.TypeInt,
									},
									"timeout": {
										Description: "Seconds to wait for activity before stopping and retrying the connection",
										Optional:    true,
										Type:        schema.TypeInt,
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
			"maven": {
				Description: "Maven specific configuration of the repository",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_policy": {
							Description: "What type of artifacts does this repository store",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"layout_policy": {
							Description: "Validate that all paths are maven artifact or metadata paths",
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
							Required:    true,
						},
					},
				},
			},
			"storage": {
				Description: "The storage configuration of the repository",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_store_name": {
							Type:        schema.TypeString,
							Default:     "default",
							Optional:    true,
							Description: "Blob store used to store repository contents",
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
							}, true),
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceRepositoryRead(d, m)
}
