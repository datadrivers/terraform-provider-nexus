package other

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/capability"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCapability() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage Nexus capabilities.",

		Create: resourceCapabilityCreate,
		Read:   resourceCapabilityRead,
		Update: resourceCapabilityUpdate,
		Delete: resourceCapabilityDelete,
		Exists: resourceCapabilityExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"type": {
				Description: "The type of capability (e.g., 'OutreachManagementCapability', 'HttpClientCapability').",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Changing type requires recreation
			},
			"notes": {
				Description: "Free-form notes about the capability.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the capability is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"properties": {
				Description: "Type-specific configuration properties.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getCapabilityFromResourceData(d *schema.ResourceData) capability.CapabilityCreate {
	properties := make(map[string]string)
	if props, ok := d.GetOk("properties"); ok {
		for k, v := range props.(map[string]interface{}) {
			properties[k] = v.(string)
		}
	}

	return capability.CapabilityCreate{
		Type:       d.Get("type").(string),
		Notes:      d.Get("notes").(string),
		Enabled:    d.Get("enabled").(bool),
		Properties: properties,
	}
}

func getCapabilityUpdateFromResourceData(d *schema.ResourceData) capability.CapabilityUpdate {
	properties := make(map[string]string)
	if props, ok := d.GetOk("properties"); ok {
		for k, v := range props.(map[string]interface{}) {
			properties[k] = v.(string)
		}
	}

	return capability.CapabilityUpdate{
		ID:         d.Id(),
		Type:       d.Get("type").(string),
		Notes:      d.Get("notes").(string),
		Enabled:    d.Get("enabled").(bool),
		Properties: properties,
	}
}

func setCapabilityToResourceData(cap *capability.Capability, d *schema.ResourceData) error {
	d.SetId(cap.ID)
	d.Set("type", cap.Type)
	d.Set("notes", cap.Notes)
	d.Set("enabled", cap.Enabled)
	d.Set("properties", cap.Properties)
	return nil
}

func resourceCapabilityCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	cap := getCapabilityFromResourceData(d)

	created, err := client.Capability.Create(cap)
	if err != nil {
		return err
	}

	d.SetId(created.ID)
	return resourceCapabilityRead(d, m)
}

func resourceCapabilityRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	cap, err := client.Capability.Get(d.Id())
	if err != nil {
		return err
	}

	if cap == nil {
		d.SetId("")
		return nil
	}

	return setCapabilityToResourceData(cap, d)
}

func resourceCapabilityUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	cap := getCapabilityUpdateFromResourceData(d)

	if err := client.Capability.Update(d.Id(), cap); err != nil {
		return err
	}

	return resourceCapabilityRead(d, m)
}

func resourceCapabilityDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Capability.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceCapabilityExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	cap, err := client.Capability.Get(d.Id())
	return cap != nil, err
}
