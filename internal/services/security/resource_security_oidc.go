package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityOIDC() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature (Nexus Repository Pro 3.93.0 or later, with ` + "`nexus.security.oauth2.enabled=true`" + ` and ` + "`nexus.jwt.enabled=true`" + ` in ` + "`nexus.properties`" + `)

Use this resource to configure Nexus Security OpenID Connect (OAuth2) integration.

The ` + "`OAuth2 Realm`" + ` must be active (via ` + "`nexus_security_realms`" + `) for external login to succeed.`,

		Create: resourceSecurityOIDCUpdate,
		Read:   resourceSecurityOIDCRead,
		Update: resourceSecurityOIDCUpdate,
		Delete: resourceSecurityOIDCDelete,
		Exists: resourceSecurityOIDCExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"client_id": {
				Description: "Unique identifier (client ID) issued by the OpenID Provider.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"client_secret": {
				Description: "Client secret issued by the OpenID Provider.",
				Required:    true,
				Type:        schema.TypeString,
				Sensitive:   true,
			},
			"authorization_url": {
				Description: "Authorization endpoint URL of the OpenID Provider (`idpAuthorizationUrl`).",
				Required:    true,
				Type:        schema.TypeString,
			},
			"token_url": {
				Description: "Token endpoint URL of the OpenID Provider (`idpTokenUrl`).",
				Required:    true,
				Type:        schema.TypeString,
			},
			"jwks_url": {
				Description: "JSON Web Key Set endpoint URL of the OpenID Provider (`idpJwksUrl`).",
				Required:    true,
				Type:        schema.TypeString,
			},
			"jws_algorithm": {
				Description: "JWT signature algorithm advertised by the OpenID Provider (`idpJwsAlgorithm`), for example `RS256`.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"username_claim": {
				Description: "ID token claim that uniquely identifies the user.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"groups_claim": {
				Description: "ID token claim carrying the user's group memberships (required for role mapping).",
				Required:    true,
				Type:        schema.TypeString,
			},
			"logout_url": {
				Description: "Logout (end-session) endpoint URL of the OpenID Provider (`idpLogoutUrl`).",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"jwks": {
				Description: "Inline JWKS JSON content (`idpJwks`). Use only when `jwks_url` cannot be reached and Nexus must validate tokens with a static key set.",
				Optional:    true,
				Type:        schema.TypeString,
				Sensitive:   true,
			},
			"first_name_claim": {
				Description: "ID token claim mapped to the user's given name.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"last_name_claim": {
				Description: "ID token claim mapped to the user's surname.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"email_claim": {
				Description: "ID token claim mapped to the user's email address.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"use_trust_store": {
				Description: "Validate the OpenID Provider certificate against the Nexus Repository truststore (`useTrustStore`).",
				Optional:    true,
				Default:     false,
				Type:        schema.TypeBool,
			},
			"exact_match_claims": {
				Description: "Claims that must match exactly for a token to be accepted (`exactMatchClaims`).",
				Optional:    true,
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"authorization_custom_params": {
				Description: "Extra query parameters appended to the authorization request (`authorizationCustomParams`).",
				Optional:    true,
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"token_request_custom_params": {
				Description: "Extra parameters appended to the token request (`tokenRequestCustomParams`).",
				Optional:    true,
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityOIDCRead(d *schema.ResourceData, m interface{}) error {
	_ = m.(*nexus.NexusClient)

	svc, err := oidc()
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

	return setSecurityOIDCToResourceData(cfg, d)
}

func resourceSecurityOIDCUpdate(d *schema.ResourceData, m interface{}) error {
	_ = m.(*nexus.NexusClient)

	svc, err := oidc()
	if err != nil {
		return err
	}

	cfg := getSecurityOIDCFromResourceData(d)
	if err := svc.Apply(cfg); err != nil {
		return err
	}
	if err := setSecurityOIDCToResourceData(&cfg, d); err != nil {
		return err
	}
	return resourceSecurityOIDCRead(d, m)
}

func resourceSecurityOIDCDelete(d *schema.ResourceData, m interface{}) error {
	_ = m.(*nexus.NexusClient)

	svc, err := oidc()
	if err != nil {
		return err
	}
	return svc.Delete()
}

func resourceSecurityOIDCExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_ = m.(*nexus.NexusClient)

	svc, err := oidc()
	if err != nil {
		return false, err
	}
	cfg, _ := svc.Read()
	// The endpoint always returns a payload (with empty strings when unset);
	// treat empty client_id as "not configured".
	return cfg != nil && cfg.ClientID != "", nil
}

// clientSecretPlaceholder is returned by Nexus when reading back the OIDC
// configuration; the real secret is never exposed via the API.
const clientSecretPlaceholder = "#~NXRM~PLACEHOLDER~PASSWORD~#"

func setSecurityOIDCToResourceData(o *OIDC, d *schema.ResourceData) error {
	d.SetId("oidc")
	d.Set("client_id", o.ClientID)
	if o.ClientSecret != "" && o.ClientSecret != clientSecretPlaceholder {
		d.Set("client_secret", o.ClientSecret)
	}
	d.Set("authorization_url", o.IdpAuthorizationURL)
	d.Set("token_url", o.IdpTokenURL)
	d.Set("jwks_url", o.IdpJwksURL)
	d.Set("jws_algorithm", o.IdpJwsAlgorithm)
	d.Set("username_claim", o.UsernameClaim)
	d.Set("groups_claim", o.GroupsClaim)
	d.Set("logout_url", o.IdpLogoutURL)
	d.Set("jwks", o.IdpJwks)
	d.Set("first_name_claim", o.FirstNameClaim)
	d.Set("last_name_claim", o.LastNameClaim)
	d.Set("email_claim", o.EmailClaim)
	d.Set("use_trust_store", o.UseTrustStore)
	d.Set("exact_match_claims", o.ExactMatchClaims)
	d.Set("authorization_custom_params", o.AuthorizationCustomParams)
	d.Set("token_request_custom_params", o.TokenRequestCustomParams)
	return nil
}

func getSecurityOIDCFromResourceData(d *schema.ResourceData) OIDC {
	o := OIDC{
		ClientID:            d.Get("client_id").(string),
		ClientSecret:        d.Get("client_secret").(string),
		IdpAuthorizationURL: d.Get("authorization_url").(string),
		IdpTokenURL:         d.Get("token_url").(string),
		IdpJwksURL:          d.Get("jwks_url").(string),
		IdpJwsAlgorithm:     d.Get("jws_algorithm").(string),
		UsernameClaim:       d.Get("username_claim").(string),
		GroupsClaim:         d.Get("groups_claim").(string),
		IdpLogoutURL:        d.Get("logout_url").(string),
		IdpJwks:             d.Get("jwks").(string),
		FirstNameClaim:      d.Get("first_name_claim").(string),
		LastNameClaim:       d.Get("last_name_claim").(string),
		EmailClaim:          d.Get("email_claim").(string),
		UseTrustStore:       d.Get("use_trust_store").(bool),
	}
	o.ExactMatchClaims = stringMapFromResource(d, "exact_match_claims")
	o.AuthorizationCustomParams = stringMapFromResource(d, "authorization_custom_params")
	o.TokenRequestCustomParams = stringMapFromResource(d, "token_request_custom_params")
	return o
}

func stringMapFromResource(d *schema.ResourceData, key string) map[string]string {
	raw, ok := d.GetOk(key)
	out := map[string]string{}
	if !ok {
		return out
	}
	for k, v := range raw.(map[string]interface{}) {
		out[k] = v.(string)
	}
	return out
}
