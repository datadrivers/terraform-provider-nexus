package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceScriptCreate,
		Read:   resourceScriptRead,
		Update: resourceScriptUpdate,
		Delete: resourceScriptDelete,
		Exists: resourceScriptExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the script.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"content": {
				Description: "The content of this script.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"type": {
				Description: "The type of the script. Default: `groovy`",
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "groovy",
			},
		},
	}
}

func getScriptFromResourceData(d *schema.ResourceData) nexus.Script {
	return nexus.Script{
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
		Type:    d.Get("type").(string),
	}
}

func resourceScriptCreate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	script := getScriptFromResourceData(d)

	if err := nexusClient.ScriptCreate(&script); err != nil {
		return err
	}

	if err := nexusClient.ScriptRun(script.Name); err != nil {
		return err
	}

	d.SetId(script.Name)
	return resourceScriptRead(d, m)
}

func resourceScriptRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	script, err := nexusClient.ScriptRead(d.Id())
	if err != nil {
		return err
	}

	if script == nil {
		d.SetId("")
		return nil
	}

	d.SetId(script.Name)
	d.Set("name", script.Name)
	d.Set("type", script.Type)
	d.Set("content", script.Content)

	return nil
}

func resourceScriptUpdate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if d.HasChange("content") || d.HasChange("type") {
		script := getScriptFromResourceData(d)
		if err := nexusClient.ScriptUpdate(&script); err != nil {
			return err
		}

		if err := nexusClient.ScriptRun(script.Name); err != nil {
			return err
		}
	}

	return resourceScriptRead(d, m)
}

func resourceScriptDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if err := nexusClient.ScriptDelete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceScriptExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	script, err := nexusClient.ScriptRead(d.Id())
	if err != nil {
		return false, err
	}
	return script != nil, nil
}
