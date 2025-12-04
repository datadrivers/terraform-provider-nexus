package other

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/iq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIQServer() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to configure the IQ Server connection in Nexus Repository Manager.",

		Create: resourceIQServerCreate,
		Read:   resourceIQServerRead,
		Update: resourceIQServerUpdate,
		Delete: resourceIQServerDelete,
		Exists: resourceIQServerExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Resource ID (always 'iqserver')",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Whether the IQ Server integration is enabled",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"show_link": {
				Description: "Whether to show the IQ Server link in the UI",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"url": {
				Description: "URL of the IQ Server",
				Type:        schema.TypeString,
				Required:    true,
			},
			"authentication_type": {
				Description: "Authentication type (USER or PKI)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": {
				Description: "Username for IQ Server authentication (required if authentication_type is USER)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "Password for IQ Server authentication (required if authentication_type is USER)",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"use_trust_store_for_url": {
				Description: "Use trust store for URL validation",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"timeout_seconds": {
				Description: "Timeout in seconds for IQ Server requests",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
			},
			"properties": {
				Description: "Additional properties",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"fail_open_mode_enabled": {
				Description: "Whether to enable fail-open mode",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func getIQServerConfigFromResourceData(d *schema.ResourceData) iq.IQServerConfiguration {
	config := iq.IQServerConfiguration{
		Enabled:             d.Get("enabled").(bool),
		ShowLink:            d.Get("show_link").(bool),
		UseTrustStoreForURL: d.Get("use_trust_store_for_url").(bool),
		FailOpenModeEnabled: d.Get("fail_open_mode_enabled").(bool),
	}

	if v, ok := d.GetOk("url"); ok {
		val := v.(string)
		config.URL = &val
	}

	if v, ok := d.GetOk("authentication_type"); ok {
		val := v.(string)
		config.AuthenticationType = &val
	}

	if v, ok := d.GetOk("username"); ok {
		val := v.(string)
		config.Username = &val
	}

	if v, ok := d.GetOk("password"); ok {
		val := v.(string)
		config.Password = &val
	}

	if v, ok := d.GetOk("timeout_seconds"); ok {
		val := v.(int)
		config.TimeoutSeconds = &val
	}

	if v, ok := d.GetOk("properties"); ok {
		val := v.(string)
		config.Properties = &val
	}

	return config
}

func setIQServerResourceData(d *schema.ResourceData, config *iq.IQServerConfiguration) error {
	d.SetId("iqserver")
	d.Set("enabled", config.Enabled)
	d.Set("show_link", config.ShowLink)
	d.Set("use_trust_store_for_url", config.UseTrustStoreForURL)
	d.Set("fail_open_mode_enabled", config.FailOpenModeEnabled)

	if config.URL != nil {
		d.Set("url", *config.URL)
	}

	if config.AuthenticationType != nil {
		d.Set("authentication_type", *config.AuthenticationType)
	}

	if config.Username != nil {
		d.Set("username", *config.Username)
	}

	// Password is not returned by the API for security reasons
	// Keep the existing value from the state

	if config.TimeoutSeconds != nil {
		d.Set("timeout_seconds", *config.TimeoutSeconds)
	}

	if config.Properties != nil {
		d.Set("properties", *config.Properties)
	}

	return nil
}

func resourceIQServerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// IQ Server configuration always exists (it's a singleton)
	return true, nil
}

func resourceIQServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	config := getIQServerConfigFromResourceData(d)

	if err := client.IQServer.Update(config); err != nil {
		return err
	}

	d.SetId("iqserver")

	return resourceIQServerRead(d, m)
}

func resourceIQServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	config, err := client.IQServer.Get()
	if err != nil {
		return err
	}

	if config == nil {
		d.SetId("")
		return nil
	}

	return setIQServerResourceData(d, config)
}

func resourceIQServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	config := getIQServerConfigFromResourceData(d)

	if err := client.IQServer.Update(config); err != nil {
		return err
	}

	return resourceIQServerRead(d, m)
}

func resourceIQServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	// Disable IQ Server integration by setting enabled to false
	// We keep the configuration but disable it
	config, err := client.IQServer.Get()
	if err != nil {
		return err
	}

	config.Enabled = false
	if err := client.IQServer.Update(*config); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
