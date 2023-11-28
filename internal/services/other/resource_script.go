package other

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	nexusSchema "github.com/dre2004/go-nexus-client/nexus3/schema"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceScript() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create and execute a custom script.",

		Create: resourceScriptCreate,
		Read:   resourceScriptRead,
		Update: resourceScriptUpdate,
		Delete: resourceScriptDelete,
		Exists: resourceScriptExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
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

func getScriptFromResourceData(d *schema.ResourceData) nexusSchema.Script {
	return nexusSchema.Script{
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
		Type:    d.Get("type").(string),
	}
}

func resourceScriptCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	script := getScriptFromResourceData(d)

	if err := client.Script.Create(&script); err != nil {
		return err
	}
	// TODO: It should be possible to configure whether to run script or not
	if err := client.Script.Run(script.Name); err != nil {
		return err
	}

	d.SetId(script.Name)
	return resourceScriptRead(d, m)
}

func resourceScriptRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	script, err := client.Script.Get(d.Id())
	if err != nil {
		return err
	}

	if script == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", script.Name)
	d.Set("type", script.Type)
	d.Set("content", script.Content)

	return nil
}

func resourceScriptUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if d.HasChange("content") || d.HasChange("type") {
		script := getScriptFromResourceData(d)
		if err := client.Script.Update(&script); err != nil {
			return err
		}

		if err := client.Script.Run(script.Name); err != nil {
			return err
		}
	}

	return resourceScriptRead(d, m)
}

func resourceScriptDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Script.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceScriptExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	script, err := client.Script.Get(d.Id())
	return script != nil, err
}
