package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityAnonymous() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to change the anonymous configuration of the nexus repository manager.",

		Create: resourceSecurityAnonymousUpdate,
		Read:   resourceSecurityAnonymousRead,
		Update: resourceSecurityAnonymousUpdate,
		Delete: resourceSecurityAnonymousDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Activate the anonymous access to the repository manager. Default: false",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"user_id": {
				Description: "The user id used by anonymous access. Default: \"anonymous\"",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "anonymous",
			},
			"realm_name": {
				Description: "The name of the used realm. Default: \"NexusAuthorizingRealm\"",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NexusAuthorizingRealm",
			},
		},
	}
}

func getAnonymousFromResourceData(d *schema.ResourceData) security.AnonymousAccessSettings {
	return security.AnonymousAccessSettings{
		Enabled:   d.Get("enabled").(bool),
		UserID:    d.Get("user_id").(string),
		RealmName: d.Get("realm_name").(string),
	}
}

func setAnonymousToResourceData(anonymous *security.AnonymousAccessSettings, d *schema.ResourceData) error {
	d.SetId("anonymous")
	d.Set("enabled", anonymous.Enabled)
	d.Set("user_id", anonymous.UserID)
	d.Set("realm_name", anonymous.RealmName)
	return nil
}

func resourceSecurityAnonymousRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	anonymous, err := client.Security.Anonymous.Read()
	if err != nil {
		return err
	}

	return setAnonymousToResourceData(anonymous, d)
}

func resourceSecurityAnonymousUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	anonymous := getAnonymousFromResourceData(d)
	if err := client.Security.Anonymous.Update(anonymous); err != nil {
		return err
	}

	return resourceSecurityAnonymousRead(d, m)
}

func resourceSecurityAnonymousDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
