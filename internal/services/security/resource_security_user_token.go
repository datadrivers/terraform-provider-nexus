package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityUserToken() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this resource to manage the global configuration for the user-tokens.`,

		Create: resourceSecurityUserTokenUpdate,
		Read:   resourceSecurityUserTokenRead,
		Update: resourceSecurityUserTokenUpdate,
		Delete: resourceSecurityUserTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"enabled": {
				Description: "Activate the feature of user tokens.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"protect_content": {
				Description: "Require user tokens for repository authentication. This does not effect UI access.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"expiration_enabled": {
				Description: "Set user tokens expiration.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"expiration_days": {
				Description: "Number of days for which you want user tokens to remain valid.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
			},
		},
	}
}

func getSecurityUserTokenFromResourceData(d *schema.ResourceData) security.UserTokenConfiguration {
	return security.UserTokenConfiguration{
		Enabled:           d.Get("enabled").(bool),
		ProtectContent:    d.Get("protect_content").(bool),
		ExpirationEnabled: d.Get("expiration_enabled").(bool),
		ExpirationDays:    d.Get("expiration_days").(int),
	}
}

func setSecurityUserTokenToResourceData(token *security.UserTokenConfiguration, d *schema.ResourceData) {
	d.SetId("golbalUserTokenConfiguration")
	d.Set("enabled", token.Enabled)
	d.Set("protect_content", token.ProtectContent)
	d.Set("expiration_enabled", token.ExpirationEnabled)
	d.Set("expiration_days", token.ExpirationDays)
}

func resourceSecurityUserTokenRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	token, err := client.Security.UserTokens.Get()
	if err != nil {
		return err
	}
	setSecurityUserTokenToResourceData(token, d)
	return nil
}

func resourceSecurityUserTokenUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	token := getSecurityUserTokenFromResourceData(d)
	if err := client.Security.UserTokens.Configure(token); err != nil {
		return err
	}

	return resourceSecurityUserTokenRead(d, m)
}

func resourceSecurityUserTokenDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
