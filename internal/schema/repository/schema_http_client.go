package repository

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	ResourceHTTPClient = &schema.Schema{
		Description: "HTTP Client configuration for proxy repositories",
		Required:    true,
		MaxItems:    1,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"authentication": ResourceHTTPClientAuthentication,
				"auto_block": {
					Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"blocked": {
					Description: "Whether to block outbound connections on the repository",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"connection": ResourceHTTPClientConnection,
			},
		},
	}

	ResourceHTTPClientWithPreemptiveAuth = &schema.Schema{
		Description: "HTTP Client configuration for proxy repositories",
		Required:    true,
		MaxItems:    1,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"authentication": ResourceHTTPClientAuthenticationWithPreemptive,
				"auto_block": {
					Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"blocked": {
					Description: "Whether to block outbound connections on the repository",
					Required:    true,
					Type:        schema.TypeBool,
				},
				"connection": ResourceHTTPClientConnection,
			},
		},
	}

	ResourceHTTPClientConnection = &schema.Schema{
		Description: "Connection configuration of the HTTP client",
		Computed:    true,
		Optional:    true,
		MaxItems:    1,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enable_circular_redirects": {
					Description: "Whether to enable redirects to the same location (may be required by some servers)",
					Optional:    true,
					Type:        schema.TypeBool,
					Default:     false,
				},
				"enable_cookies": {
					Description: "Whether to allow cookies to be stored and used",
					Optional:    true,
					Type:        schema.TypeBool,
					Default:     false,
				},
				"retries": {
					Description:      "Total retries if the initial connection attempt suffers a timeout",
					Optional:         true,
					Type:             schema.TypeInt,
					Default:          0,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 10)),
				},
				"timeout": {
					Description: "Seconds to wait for activity before stopping and retrying the connection",
					Optional:    true,
					Type:        schema.TypeInt,
					ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
						maxTimeoutSeconds := 3600
						minTimeoutSeconds := 1
						var diags diag.Diagnostics

						if v != nil {
							timeout := v.(int)
							if timeout > maxTimeoutSeconds && timeout < minTimeoutSeconds {
								diag := diag.Diagnostic{
									Severity:      diag.Error,
									Summary:       "The timeout value is wrong",
									Detail:        fmt.Sprintf("%q is not between %q and %q", timeout, minTimeoutSeconds, maxTimeoutSeconds),
									AttributePath: p,
								}
								diags = append(diags, diag)
							}
						}
						return diags
					},
					// ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 3600)),
				},
				"user_agent_suffix": {
					Description: "Custom fragment to append to User-Agent header in HTTP requests",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"use_trust_store": {
					Description: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
					Optional:    true,
					Type:        schema.TypeBool,
					Default:     false,
				},
			},
		},
	}

	ResourceHTTPClientAuthentication = &schema.Schema{
		Description: "Authentication configuration of the HTTP client",
		MaxItems:    1,
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description:      "Authentication type. Possible values: `ntlm` or `username`",
					Required:         true,
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ntlm", "username"}, false)),
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
	}

	ResourceHTTPClientAuthenticationWithPreemptive = &schema.Schema{
		Description: "Authentication configuration of the HTTP client",
		MaxItems:    1,
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description:  "Authentication type. Possible values: `ntlm` or `username`",
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ntlm", "username"}, false),
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
				"preemptive": {
					Description: "Whether to use pre-emptive authentication. Use with caution. Defaults to false.",
					Optional:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}

	DataSourceHTTPClient = &schema.Schema{
		Description: "HTTP Client configuration for proxy repositories",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"authentication": DataSourceHTTPClientAuthentication,
				"auto_block": {
					Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"blocked": {
					Description: "Whether to block outbound connections on the repository",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"connection": DataSourceHTTPClientConnection,
			},
		},
	}

	DataSourceHTTPClientWithPreemptiveAuth = &schema.Schema{
		Description: "HTTP Client configuration for proxy repositories",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"authentication": DataSourceHTTPClientAuthenticationWithPreemptive,
				"auto_block": {
					Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"blocked": {
					Description: "Whether to block outbound connections on the repository",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"connection": DataSourceHTTPClientConnection,
			},
		},
	}

	DataSourceHTTPClientConnection = &schema.Schema{
		Description: "Connection configuration of the HTTP client",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enable_circular_redirects": {
					Description: "Whether to enable redirects to the same location (may be required by some servers)",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"enable_cookies": {
					Description: "Whether to allow cookies to be stored and used",
					Computed:    true,
					Type:        schema.TypeBool,
				},
				"retries": {
					Description: "Total retries if the initial connection attempt suffers a timeout",
					Computed:    true,
					Type:        schema.TypeInt,
				},
				"timeout": {
					Description: "Seconds to wait for activity before stopping and retrying the connection",
					Computed:    true,
					Type:        schema.TypeInt,
				},
				"user_agent_suffix": {
					Description: "Custom fragment to append to User-Agent header in HTTP requests",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"use_trust_store": {
					Description: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
					Computed:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}

	DataSourceHTTPClientAuthentication = &schema.Schema{
		Description: "Authentication configuration of the HTTP client",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "Authentication type. Possible values: `ntlm` or `username`",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"username": {
					Description: "The username used by the proxy repository",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"password": {
					Description: "The password used by the proxy repository",
					Computed:    true,
					Sensitive:   true,
					Type:        schema.TypeString,
				},
				"ntlm_domain": {
					Description: "The ntlm domain to connect",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"ntlm_host": {
					Description: "The ntlm host to connect",
					Computed:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	DataSourceHTTPClientAuthenticationWithPreemptive = &schema.Schema{
		Description: "Authentication configuration of the HTTP client",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "Authentication type. Possible values: `ntlm` or `username`",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"username": {
					Description: "The username used by the proxy repository",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"password": {
					Description: "The password used by the proxy repository",
					Computed:    true,
					Sensitive:   true,
					Type:        schema.TypeString,
				},
				"ntlm_domain": {
					Description: "The ntlm domain to connect",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"ntlm_host": {
					Description: "The ntlm host to connect",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"preemptive": {
					Description: "Whether to use pre-emptive authentication",
					Computed:    true,
					Type:        schema.TypeBool,
				},
			},
		},
	}
)
