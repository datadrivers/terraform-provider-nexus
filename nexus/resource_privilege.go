package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourcePrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivilegeCreate,
		Read:   resourcePrivilegeRead,
		Update: resourcePrivilegeUpdate,
		Delete: resourcePrivilegeDelete,
		Exists: resourcePrivilegeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"actions": {
				Description: "Actions for the privilege (browse, read, edit, add, delete, all and run)",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Type:        schema.TypeSet,
			},
			"content_selector": {
				Description: "The content selector for the privilege",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A description of the privilege",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"domain": {
				Description: "The domain of the privilege",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"format": {
				Description:  "The format of the privilege",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(nexus.RepositoryFormats, false),
			},
			"name": {
				Description: "The name of the privilege",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			"repository": {
				Description: "The repository of the privilege",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"type": {
				Description:  "The type of the privilege",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"application", "repository-admin", "repository-content-selector", "repository-view", "script", "wildcard"}, false),
			},
		},
	}
}

func getPrivilegeFromResourceData(d *schema.ResourceData) nexus.Privilege {
	privilege := nexus.Privilege{
		Actions: interfaceSliceToStringSlice(d.Get("actions").(*schema.Set).List()),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		privilege.Description = description.(string)
	}

	if contentSelector, ok := d.GetOk("content_selector"); ok {
		privilege.ContentSelector = contentSelector.(string)
	}

	if domain, ok := d.GetOk("domain"); ok {
		privilege.Domain = domain.(string)
	}

	if format, ok := d.GetOk("format"); ok {
		privilege.Format = format.(string)
	}

	if repository, ok := d.GetOk("repository"); ok {
		privilege.Repository = repository.(string)
	}

	return privilege
}

func setPrivilegeToResourceData(privilege *nexus.Privilege, d *schema.ResourceData) error {
	d.SetId(privilege.Name)
	d.Set("actions", privilege.Actions)
	d.Set("content_selector", privilege.ContentSelector)
	d.Set("description", privilege.Description)
	d.Set("domain", privilege.Domain)
	d.Set("format", privilege.Format)
	d.Set("name", privilege.Name)
	d.Set("repository", privilege.Repository)
	d.Set("type", privilege.Type)
	return nil
}

func resourcePrivilegeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	privilege := getPrivilegeFromResourceData(d)

	if err := client.PrivilegeCreate(privilege); err != nil {
		return err
	}

	d.SetId(privilege.Name)

	return resourcePrivilegeRead(d, m)
}

func resourcePrivilegeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	privilege, err := client.PrivilegeRead(d.Id())
	if err != nil {
		return err
	}

	if privilege == nil {
		d.SetId("")
		return nil
	}

	return setPrivilegeToResourceData(privilege, d)
}

func resourcePrivilegeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	privilege := getPrivilegeFromResourceData(d)
	if err := client.PrivilegeUpdate(d.Id(), privilege); err != nil {
		return err
	}

	return resourcePrivilegeRead(d, m)
}

func resourcePrivilegeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)

	if err := client.PrivilegeDelete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourcePrivilegeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(nexus.Client)

	privilege, err := client.PrivilegeRead(d.Id())
	return privilege != nil, err
}
