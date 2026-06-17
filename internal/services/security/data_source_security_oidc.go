package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityOIDC() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature (Nexus Repository Pro 3.93.0 or later)

Use this data source to read the Nexus OAuth2/OpenID Connect configuration.`,

		Read: dataSourceSecurityOIDCRead,
		Schema: map[string]*schema.Schema{
			"id":                common.DataSourceID,
			"client_id":         {Computed: true, Type: schema.TypeString, Description: "OIDC client identifier."},
			"client_secret":     {Computed: true, Type: schema.TypeString, Sensitive: true, Description: "OIDC client secret (may be redacted)."},
			"authorization_url": {Computed: true, Type: schema.TypeString, Description: "OIDC authorization endpoint URL."},
			"token_url":         {Computed: true, Type: schema.TypeString, Description: "OIDC token endpoint URL."},
			"jwks_url":          {Computed: true, Type: schema.TypeString, Description: "OIDC JWKS endpoint URL."},
			"jws_algorithm":     {Computed: true, Type: schema.TypeString, Description: "JWT signature algorithm."},
			"username_claim":    {Computed: true, Type: schema.TypeString, Description: "Username claim."},
			"groups_claim":      {Computed: true, Type: schema.TypeString, Description: "Groups claim."},
			"logout_url":        {Computed: true, Type: schema.TypeString, Description: "OIDC logout endpoint URL."},
			"jwks":              {Computed: true, Type: schema.TypeString, Sensitive: true, Description: "Inline JWKS JSON."},
			"first_name_claim":  {Computed: true, Type: schema.TypeString, Description: "Given-name claim."},
			"last_name_claim":   {Computed: true, Type: schema.TypeString, Description: "Surname claim."},
			"email_claim":       {Computed: true, Type: schema.TypeString, Description: "Email claim."},
			"use_trust_store":   {Computed: true, Type: schema.TypeBool, Description: "Whether the Nexus truststore is used."},
			"exact_match_claims": {
				Computed: true, Type: schema.TypeMap,
				Elem: &schema.Schema{Type: schema.TypeString}, Description: "Claims that must match exactly.",
			},
			"authorization_custom_params": {
				Computed: true, Type: schema.TypeMap,
				Elem: &schema.Schema{Type: schema.TypeString}, Description: "Extra authorization-request parameters.",
			},
			"token_request_custom_params": {
				Computed: true, Type: schema.TypeMap,
				Elem: &schema.Schema{Type: schema.TypeString}, Description: "Extra token-request parameters.",
			},
		},
	}
}

func dataSourceSecurityOIDCRead(d *schema.ResourceData, m interface{}) error {
	return resourceSecurityOIDCRead(d, m)
}
