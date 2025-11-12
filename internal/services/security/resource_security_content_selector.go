package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityContentSelector() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus Content Selector.",

		Create: resourceSecurityContentSelectorCreate,
		Read:   resourceSecurityContentSelectorRead,
		Update: resourceSecurityContentSelectorUpdate,
		Delete: resourceSecurityContentSelectorDelete,
		Exists: resourceSecurityContentSelectorExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"name": {
				Description: "Content selector name",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A description of the content selector",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"expression": {
				Description: "The content selector expression",
				Required:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func getContentSelectorFromResourceData(d *schema.ResourceData) security.ContentSelector {
	contentSelector := security.ContentSelector{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		contentSelector.Description = description.(string)
	}

	return contentSelector
}

func setContentSelectorToResourceData(contentSelector *security.ContentSelector, d *schema.ResourceData) error {
	d.SetId(contentSelector.Name)
	d.Set("description", contentSelector.Description)
	d.Set("expression", contentSelector.Expression)
	d.Set("name", contentSelector.Name)
	return nil
}

func resourceSecurityContentSelectorCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	contentSelector := getContentSelectorFromResourceData(d)

	if err := client.Security.ContentSelector.Create(contentSelector); err != nil {
		return err
	}

	d.SetId(contentSelector.Name)

	return resourceSecurityContentSelectorRead(d, m)
}

func resourceSecurityContentSelectorRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	contentSelector, err := client.Security.ContentSelector.Get(d.Id())
	if err != nil {
		return err
	}

	if contentSelector == nil {
		d.SetId("")
		return nil
	}

	return setContentSelectorToResourceData(contentSelector, d)
}

func resourceSecurityContentSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	contentSelector := getContentSelectorFromResourceData(d)
	if err := client.Security.ContentSelector.Update(d.Id(), contentSelector); err != nil {
		return err
	}

	return resourceSecurityContentSelectorRead(d, m)
}

func resourceSecurityContentSelectorDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.ContentSelector.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceSecurityContentSelectorExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	contentSelector, err := client.Security.ContentSelector.Get(d.Id())
	return contentSelector != nil, err
}
