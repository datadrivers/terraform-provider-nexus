package security

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
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
		},
	}
}

func getSecurityUserTokenFromResourceData(d *schema.ResourceData) security.UserTokenConfiguration {
	return security.UserTokenConfiguration{
		Enabled:        d.Get("enabled").(bool),
		ProtectContent: d.Get("protect_content").(bool),
	}
}

func setSecurityUserTokenToResourceData(token *security.UserTokenConfiguration, d *schema.ResourceData) {
	d.SetId("golbalUserTokenConfiguration")
	d.Set("enabled", token.Enabled)
	d.Set("protect_content", token.ProtectContent)
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
