package repository

import (
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getResourceHTTPClientSchema() *schema.Schema {
	return &schema.Schema{
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
	}
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
