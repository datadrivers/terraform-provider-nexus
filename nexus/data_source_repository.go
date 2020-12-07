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
				Optional: true,
				Type:     schema.TypeString,
			},
			"online": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cleanup": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_names": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
					},
				},
			},
			"apt": {
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
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
			"docker": {
				Type:     schema.TypeList,
				Optional: true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Authentication type",
										Required:    true,
										Type:        schema.TypeString,
									},
									"username": {
										Description: "",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"ntlm_domain": {
										Description: "",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"ntlm_host": {
										Description: "",
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
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_policy": {
							Default:  "RELEASE",
							Type:     schema.TypeString,
							Optional: true,
						},
						"layout_policy": {
							Default:  "PERMISSIVE",
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
