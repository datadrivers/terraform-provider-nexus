package other

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const HTTPClientConfigID = "http"

func ResourceHTTPClient() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to configure Nexus system HTTP/HTTPS proxy settings (`GET/PUT/DELETE /service/rest/v1/http`). There is exactly one HTTP configuration per Nexus instance.",

		Create: resourceHTTPClientCreate,
		Read:   resourceHTTPClientRead,
		Update: resourceHTTPClientUpdate,
		Delete: resourceHTTPClientDelete,
		Exists: resourceHTTPClientExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"timeout": {
				Description:  "Connection/request timeout in seconds.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"retries": {
				Description:  "Number of retries for outbound HTTP requests.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(0, 10),
			},
			"user_agent": {
				Description: "Custom User-Agent suffix for outbound HTTP requests.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"non_proxy_hosts": {
				Description: "Host patterns that should bypass the HTTP/HTTPS proxy.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"http_proxy":  resourceHTTPProxySchema("HTTP proxy configuration."),
			"https_proxy": resourceHTTPProxySchema("HTTPS proxy configuration."),
		},
	}
}

func resourceHTTPProxySchema(desc string) *schema.Schema {
	return &schema.Schema{
		Description: desc,
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Description: "Whether this proxy is enabled.",
					Type:        schema.TypeBool,
					Required:    true,
				},
				"host": {
					Description: "Proxy hostname.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"port": {
					Description:  "Proxy port.",
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 65535),
				},
				"username": {
					Description: "Proxy authentication username.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"password": {
					Description: "Proxy authentication password.",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
				},
				"ntlm_host": {
					Description: "NTLM host for proxy authentication.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"ntlm_domain": {
					Description: "NTLM domain for proxy authentication.",
					Type:        schema.TypeString,
					Optional:    true,
				},
			},
		},
	}
}

func resourceHTTPClientRead(d *schema.ResourceData, m interface{}) error {
	nc := m.(*nexus.NexusClient)
	svc, err := httpClient(nc)
	if err != nil {
		return err
	}

	cfg, err := svc.Read()
	if err != nil {
		return err
	}
	if cfg == nil {
		d.SetId("")
		return nil
	}

	return setHTTPClientToResourceData(cfg, d)
}

func resourceHTTPClientCreate(d *schema.ResourceData, m interface{}) error {
	return resourceHTTPClientUpdate(d, m)
}

func resourceHTTPClientUpdate(d *schema.ResourceData, m interface{}) error {
	nc := m.(*nexus.NexusClient)
	svc, err := httpClient(nc)
	if err != nil {
		return err
	}

	cfg := getHTTPClientFromResourceData(d)
	if err := svc.Apply(cfg); err != nil {
		return err
	}

	d.SetId(HTTPClientConfigID)
	return resourceHTTPClientRead(d, m)
}

func resourceHTTPClientDelete(d *schema.ResourceData, m interface{}) error {
	nc := m.(*nexus.NexusClient)
	svc, err := httpClient(nc)
	if err != nil {
		return err
	}
	if err := svc.Delete(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceHTTPClientExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nc := m.(*nexus.NexusClient)
	svc, err := httpClient(nc)
	if err != nil {
		return false, err
	}
	cfg, err := svc.Read()
	if err != nil {
		return false, err
	}
	return cfg != nil, nil
}

func setHTTPClientToResourceData(cfg *HTTPClientConfig, d *schema.ResourceData) error {
	d.SetId(HTTPClientConfigID)
	d.Set("timeout", cfg.Timeout)
	d.Set("retries", cfg.Retries)
	d.Set("user_agent", cfg.UserAgent)
	d.Set("non_proxy_hosts", tools.StringSliceToInterfaceSlice(cfg.NonProxyHosts))
	if err := d.Set("http_proxy", flattenHTTPProxy(cfg.HTTPProxy, d, "http_proxy")); err != nil {
		return err
	}
	return d.Set("https_proxy", flattenHTTPProxy(cfg.HTTPSProxy, d, "https_proxy"))
}

func flattenHTTPProxy(p HTTPProxy, d *schema.ResourceData, key string) []map[string]interface{} {
	if !p.Enabled && p.Host == "" && p.portInt() == 0 {
		if _, ok := d.GetOk(key); !ok {
			return nil
		}
	}

	password := p.AuthInfo.Password
	if password == "" {
		if list, ok := d.GetOk(key); ok {
			items := list.([]interface{})
			if len(items) > 0 && items[0] != nil {
				if m, ok := items[0].(map[string]interface{}); ok {
					if v, ok := m["password"].(string); ok {
						password = v
					}
				}
			}
		}
	}

	return []map[string]interface{}{
		{
			"enabled":     p.Enabled,
			"host":        p.Host,
			"port":        p.portInt(),
			"username":    p.AuthInfo.Username,
			"password":    password,
			"ntlm_host":   p.AuthInfo.NtlmHost,
			"ntlm_domain": p.AuthInfo.NtlmDomain,
		},
	}
}

func getHTTPClientFromResourceData(d *schema.ResourceData) HTTPClientConfig {
	cfg := HTTPClientConfig{
		Timeout:       d.Get("timeout").(int),
		Retries:       d.Get("retries").(int),
		UserAgent:     d.Get("user_agent").(string),
		NonProxyHosts: []string{},
	}
	if v, ok := d.GetOk("non_proxy_hosts"); ok {
		cfg.NonProxyHosts = tools.ConvertStringSet(v.(*schema.Set))
	}
	cfg.HTTPProxy = expandHTTPProxy(d.Get("http_proxy"))
	cfg.HTTPSProxy = expandHTTPProxy(d.Get("https_proxy"))
	return cfg
}

func expandHTTPProxy(raw interface{}) HTTPProxy {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return HTTPProxy{
			Enabled:  false,
			Host:     "",
			Port:     "",
			AuthInfo: HTTPAuthInfo{},
		}
	}
	m := list[0].(map[string]interface{})
	username := m["username"].(string)
	password := m["password"].(string)
	ntlmHost := m["ntlm_host"].(string)
	ntlmDomain := m["ntlm_domain"].(string)
	authEnabled := username != "" || password != "" || ntlmHost != "" || ntlmDomain != ""

	port := m["port"].(int)
	var portVal interface{}
	if port > 0 {
		portVal = port
	} else {
		portVal = ""
	}

	return HTTPProxy{
		Enabled: m["enabled"].(bool),
		Host:    m["host"].(string),
		Port:    portVal,
		AuthInfo: HTTPAuthInfo{
			Enabled:    authEnabled,
			Username:   username,
			Password:   password,
			NtlmHost:   ntlmHost,
			NtlmDomain: ntlmDomain,
		},
	}
}
