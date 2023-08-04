package other

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema"

	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// There is exactly one mail config, so use fixed value
const MailConfigId = "cfg"

func ResourceMailConfig() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus Routing Rule.",

		Create: resourceEmailConfigCreate,
		Read:   resourceEmailConfigRead,
		Update: resourceEmailConfigUpdate,
		Delete: resourceEmailConfigDelete,

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"enabled": {
				Description: "Whether the config is enabled or not",
				Type:        schema.TypeBool,
				Computed:    true,
			},
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
		},
	}
}

func resourceEmailConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	mailconfig, err := client.MailConfig.Get()

	if err != nil {
		return err
	}

	if mailconfig == nil {
		d.SetId("")
		return nil
	}

	d.Set("enabled", mailconfig.Enabled)
	d.Set("host", mailconfig.Host)
	d.Set("fromaddress", mailconfig.FromAddress)
	d.Set("port", mailconfig.Port)

	return nil
}

func getEmailConfigFromResourceData(d *schema.ResourceData) nexusSchema.MailConfig {
	return nexusSchema.MailConfig{
		Host:        d.Get("host").(string),
		FromAddress: d.Get("from_address").(string),
		Port:        d.Get("port").(int),
	}
}

func resourceEmailConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	mailconfig := getEmailConfigFromResourceData(d)

	if err := client.MailConfig.Create(&mailconfig); err != nil {
		return err
	}

	d.SetId(MailConfigId)
	// return resourceRoutingRuleRead(d, m)
	return resourceEmailConfigRead(d, m)
}

func resourceEmailConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	mailconfig := getEmailConfigFromResourceData(d)
	if err := client.MailConfig.Update(&mailconfig); err != nil {
		return err
	}

	return resourceEmailConfigRead(d, m)
}

func resourceEmailConfigDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.MailConfig.Delete(); err != nil {
		return err
	}

	d.SetId(MailConfigId)
	return nil
}
