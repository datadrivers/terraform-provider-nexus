package other

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	nexusSchema "github.com/dre2004/go-nexus-client/nexus3/schema"

	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/dre2004/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// There is exactly one mail config, so use fixed value
const MailConfigId = "cfg"

func ResourceMailConfig() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to configure Nexus' mailing behaviour",

		Create: resourceMailConfigCreate,
		Read:   resourceMailConfigRead,
		Update: resourceMailConfigUpdate,
		Delete: resourceMailConfigDelete,
		Exists: resourceMailConfigExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"host": {
				Description: "hostname",
				Type:        schema.TypeString,
				Required:    true,
			},
			"from_address": {
				Description: "fromAddress",
				Type:        schema.TypeString,
				Required:    true,
			},
			"port": {
				Description: "port",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"enabled": {
				Description: "Whether the config is enabled or not",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"username": {
				Description: "Username",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "Password",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"subject_prefix": {
				Description: "Subject prefix",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"start_tls_enabled": {
				Description: "Star TLS Enabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"start_tls_required": {
				Description: "Star TLS required",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"ssl_on_connect_enabled": {
				Description: "SSL on connect enabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"ssl_server_identity_check_enabled": {
				Description: "SSL on connect enabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"nexus_trust_store_enabled": {
				Description: "SSL on connect enabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceMailConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	mailconfig, err := client.MailConfig.Get()

	if err != nil {
		return err
	}

	if mailconfig == nil {
		d.SetId(MailConfigId)
		return nil
	}

	d.Set("host", mailconfig.Host)
	d.Set("from_address", mailconfig.FromAddress)
	d.Set("port", mailconfig.Port)
	d.Set("enabled", mailconfig.Enabled)
	d.Set("username", mailconfig.Username)
	d.Set("subject_prefix", mailconfig.SubjectPrefix)
	d.Set("start_tls_enabled", mailconfig.StartTlsEnabled)
	d.Set("start_tls_required", mailconfig.StartTlsRequired)
	d.Set("ssl_on_connect_enabled", mailconfig.SslOnConnectEnabled)
	d.Set("ssl_server_identity_check_enabled", mailconfig.SslServerIdentityCheckEnabled)
	d.Set("nexus_trust_store_enabled", mailconfig.NexusTrustStoreEnabled)

	return nil
}

func getMailConfigFromResourceData(d *schema.ResourceData) nexusSchema.MailConfig {
	mailconfig := nexusSchema.MailConfig{
		Host:                          d.Get("host").(string),
		FromAddress:                   d.Get("from_address").(string),
		Port:                          d.Get("port").(int),
		Enabled:                       tools.GetBoolPointer(d.Get("enabled").(bool)),
		Username:                      tools.GetStringPointer(d.Get("username").(string)),
		SubjectPrefix:                 tools.GetStringPointer(d.Get("subject_prefix").(string)),
		StartTlsEnabled:               tools.GetBoolPointer(d.Get("start_tls_enabled").(bool)),
		StartTlsRequired:              tools.GetBoolPointer(d.Get("start_tls_required").(bool)),
		SslOnConnectEnabled:           tools.GetBoolPointer(d.Get("ssl_on_connect_enabled").(bool)),
		SslServerIdentityCheckEnabled: tools.GetBoolPointer(d.Get("ssl_server_identity_check_enabled").(bool)),
		NexusTrustStoreEnabled:        tools.GetBoolPointer(d.Get("nexus_trust_store_enabled").(bool)),
	}
	return mailconfig
}

func resourceMailConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	mailconfig := getMailConfigFromResourceData(d)

	if err := client.MailConfig.Create(&mailconfig); err != nil {
		return err
	}

	d.SetId(MailConfigId)
	return resourceMailConfigRead(d, m)
}

func resourceMailConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	mailconfig := getMailConfigFromResourceData(d)
	if err := client.MailConfig.Update(&mailconfig); err != nil {
		return err
	}

	return resourceMailConfigRead(d, m)
}

func resourceMailConfigDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.MailConfig.Delete(); err != nil {
		return err
	}

	d.SetId(MailConfigId)
	return nil
}

func resourceMailConfigExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	mailconfig, err := client.MailConfig.Get()
	return mailconfig != nil, err
}
