package other

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceHTTPClient() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query Nexus system HTTP/HTTPS proxy settings.",

		Read: dataSourceHTTPClientRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"timeout": {
				Description: "Connection/request timeout in seconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"retries": {
				Description: "Number of retries for outbound HTTP requests.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"user_agent": {
				Description: "Custom User-Agent suffix for outbound HTTP requests.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"non_proxy_hosts": {
				Description: "Host patterns that should bypass the HTTP/HTTPS proxy.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"http_proxy":  dataSourceHTTPProxySchema("HTTP proxy configuration."),
			"https_proxy": dataSourceHTTPProxySchema("HTTPS proxy configuration."),
		},
	}
}

func dataSourceHTTPProxySchema(desc string) *schema.Schema {
	return &schema.Schema{
		Description: desc,
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Description: "Whether this proxy is enabled.",
					Type:        schema.TypeBool,
					Computed:    true,
				},
				"host": {
					Description: "Proxy hostname.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"port": {
					Description: "Proxy port.",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				"username": {
					Description: "Proxy authentication username.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"password": {
					Description: "Proxy authentication password.",
					Type:        schema.TypeString,
					Computed:    true,
					Sensitive:   true,
				},
				"ntlm_host": {
					Description: "NTLM host for proxy authentication.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"ntlm_domain": {
					Description: "NTLM domain for proxy authentication.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

func dataSourceHTTPClientRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(HTTPClientConfigID)
	return resourceHTTPClientRead(d, m)
}
